// Package storage 负责文件上传到 MinIO/S3，含类型/大小限制与服务端 MIME 校验。
package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/Alike/backend/pkg/config"
)

// 上传相关错误。
var (
	ErrFileTooLarge      = errors.New("文件超过大小限制")
	ErrUnsupportedType   = errors.New("不支持的文件类型")
	ErrEmptyFile         = errors.New("文件为空")
	ErrStorageDisabled   = errors.New("对象存储未配置")
	sniffLen             = 512 // net/http.DetectContentType 只需前 512 字节
)

// FileCategory 是允许上传的文件大类。
type FileCategory string

const (
	CategoryImage    FileCategory = "image"
	CategoryDocument FileCategory = "document"
)

// allowedImageMIME 是允许的图片 MIME（DetectContentType 可识别）。
var allowedImageMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// allowedDocMIME 是允许的文档 MIME 与扩展名。
var allowedDocMIME = map[string]string{
	"application/pdf":    ".pdf",
	"application/zip":    ".zip",
	"text/plain":         ".txt",
	"text/plain; charset=utf-8": ".txt",
}

// UploadResult 是一次上传的返回信息。
type UploadResult struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"` // image/document
	MIME     string `json:"mime"`
	Size     int64  `json:"size"`
}

// Store 封装 MinIO 客户端与上传策略。
type Store struct {
	client       *minio.Client
	bucket       string
	publicURL    string // 对外基地址（不含 bucket），末尾无 /
	useSSL       bool
	endpoint     string
	maxImageSize int64
	maxDocSize   int64
}

// New 依据配置构造 Store。返回的 Store 可能因缺少依赖而为占位（client 为 nil），
// 调用 Upload 时会返回 ErrStorageDisabled。
func New(cfg *config.Config) (*Store, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, err
	}

	scheme := "http"
	if cfg.MinIOUseSSL {
		scheme = "https"
	}
	publicURL := strings.TrimRight(cfg.MinIOPublicURL, "/")
	if publicURL == "" {
		publicURL = scheme + "://" + cfg.MinIOEndpoint
	}

	return &Store{
		client:       client,
		bucket:       cfg.MinIOBucket,
		publicURL:    publicURL,
		useSSL:       cfg.MinIOUseSSL,
		endpoint:     cfg.MinIOEndpoint,
		maxImageSize: cfg.UploadMaxImageBytes,
		maxDocSize:   cfg.UploadMaxDocBytes,
	}, nil
}

// EnsureBucket 确保目标 bucket 存在（幂等）。DB/MinIO 不可用时由调用方容错。
func (s *Store) EnsureBucket(ctx context.Context) error {
	if s == nil || s.client == nil {
		return ErrStorageDisabled
	}
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if !exists {
		return s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
	}
	return nil
}

// classify 依据 MIME 返回文件大类、目标扩展名与大小上限。不支持则返回 ErrUnsupportedType。
func (s *Store) classify(mime string) (FileCategory, string, int64, error) {
	if ext, ok := allowedImageMIME[mime]; ok {
		return CategoryImage, ext, s.maxImageSize, nil
	}
	if ext, ok := allowedDocMIME[mime]; ok {
		return CategoryDocument, ext, s.maxDocSize, nil
	}
	return "", "", 0, ErrUnsupportedType
}

// Upload 读取文件流，做服务端 MIME 校验与大小限制后上传到 MinIO，返回可访问 URL。
// objectPrefix 用于分目录（如用户 ID）；size 为客户端声明的大小（可为 0，未知）。
func (s *Store) Upload(ctx context.Context, r io.Reader, declaredSize int64, objectPrefix string) (*UploadResult, error) {
	if s == nil || s.client == nil {
		return nil, ErrStorageDisabled
	}

	// 读取文件头做 MIME 嗅探，不信任扩展名。
	header := make([]byte, sniffLen)
	n, err := io.ReadFull(r, header)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, err
	}
	header = header[:n]
	if n == 0 {
		return nil, ErrEmptyFile
	}
	mime := http.DetectContentType(header)

	category, ext, maxSize, err := s.classify(mime)
	if err != nil {
		return nil, err
	}

	// 已声明大小时先行拦截，避免读取超大文件。
	if declaredSize > 0 && declaredSize > maxSize {
		return nil, ErrFileTooLarge
	}

	// 拼接已读文件头与剩余流，并用 LimitReader 强制大小上限（+1 字节用于探测越界）。
	full := io.MultiReader(bytes.NewReader(header), r)
	limited := io.LimitReader(full, maxSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxSize {
		return nil, ErrFileTooLarge
	}
	size := int64(len(data))

	objectName := buildObjectName(objectPrefix, string(category), ext)
	_, err = s.client.PutObject(ctx, s.bucket, objectName,
		bytes.NewReader(data), size,
		minio.PutObjectOptions{ContentType: baseMIME(mime)})
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:      s.publicURL + "/" + s.bucket + "/" + objectName,
		FileName: path.Base(objectName),
		FileType: string(category),
		MIME:     baseMIME(mime),
		Size:     size,
	}, nil
}

// baseMIME 去掉 MIME 中的参数部分（如 "; charset=utf-8"）。
func baseMIME(mime string) string {
	if i := strings.IndexByte(mime, ';'); i >= 0 {
		return strings.TrimSpace(mime[:i])
	}
	return mime
}

// buildObjectName 生成对象存储中的 key：<category>/<prefix>/<yyyymmdd>/<random><ext>。
// prefix 通常为上传者用户 ID，用于分目录；随机后缀保证唯一。
func buildObjectName(prefix, category, ext string) string {
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		prefix = "anon"
	}
	day := time.Now().UTC().Format("20060102")
	return fmt.Sprintf("%s/%s/%s/%s%s", category, prefix, day, randomHex(16), ext)
}

// randomHex 返回 n 字节的十六进制随机串；熵源失败时回退到纳秒时间戳。
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
