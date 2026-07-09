/**
 * 认证表单校验 —— 纯函数，供登录/注册页复用并可单测。
 * 返回字段级错误信息（空字符串表示无错误）。
 */
export const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function validateEmail(email: string): string {
  return EMAIL_RE.test(email.trim()) ? '' : '请输入有效的邮箱地址'
}

export function validatePassword(password: string): string {
  return password.length >= 6 ? '' : '密码至少 6 位'
}

export function validateNickname(nickname: string): string {
  return nickname.trim() ? '' : '昵称不能为空'
}

/** 工龄选填：空字符串视为未填（通过）；填写时必须为 0-60 的整数。 */
export function validateWorkYears(value: string): string {
  if (value === '') return ''
  const n = Number(value)
  if (!Number.isInteger(n) || n < 0 || n > 60) return '工龄需为 0-60 的整数'
  return ''
}

export interface RegisterFormInput {
  email: string
  password: string
  nickname: string
  work_years: string
}

export interface RegisterFieldErrors {
  [key: string]: string
  email: string
  password: string
  nickname: string
  work_years: string
}

/** 校验注册表单，返回字段错误集合。全部为空字符串即通过。 */
export function validateRegisterForm(
  form: RegisterFormInput,
): RegisterFieldErrors {
  return {
    email: validateEmail(form.email),
    password: validatePassword(form.password),
    nickname: validateNickname(form.nickname),
    work_years: validateWorkYears(form.work_years),
  }
}

/** 判断错误集合是否全部通过（无任何非空错误）。 */
export function isFormValid(errors: Record<string, string>): boolean {
  return Object.values(errors).every((e) => e === '')
}
