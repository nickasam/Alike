import { describe, it, expect } from 'vitest'
import {
  validateEmail,
  validatePassword,
  validateNickname,
  validateWorkYears,
  validateRegisterForm,
  isFormValid,
} from '~/composables/useAuthValidation'

describe('auth 表单校验', () => {
  it('邮箱格式校验', () => {
    expect(validateEmail('a@b.com')).toBe('')
    expect(validateEmail('  a@b.com  ')).toBe('') // 前后空格被 trim
    expect(validateEmail('not-an-email')).not.toBe('')
    expect(validateEmail('a@b')).not.toBe('')
    expect(validateEmail('')).not.toBe('')
  })

  it('密码至少 6 位', () => {
    expect(validatePassword('123456')).toBe('')
    expect(validatePassword('12345')).not.toBe('')
    expect(validatePassword('')).not.toBe('')
  })

  it('昵称非空（含纯空格视为空）', () => {
    expect(validateNickname('牛马')).toBe('')
    expect(validateNickname('   ')).not.toBe('')
    expect(validateNickname('')).not.toBe('')
  })

  it('工龄选填：空通过，非法值报错', () => {
    expect(validateWorkYears('')).toBe('') // 选填留空
    expect(validateWorkYears('0')).toBe('')
    expect(validateWorkYears('5')).toBe('')
    expect(validateWorkYears('60')).toBe('')
    expect(validateWorkYears('-1')).not.toBe('')
    expect(validateWorkYears('61')).not.toBe('')
    expect(validateWorkYears('3.5')).not.toBe('')
    expect(validateWorkYears('abc')).not.toBe('')
  })

  it('validateRegisterForm 通过合法表单', () => {
    const errs = validateRegisterForm({
      email: 'niuma@example.com',
      password: 'secret1',
      nickname: '牛马一号',
      work_years: '3',
    })
    expect(isFormValid(errs)).toBe(true)
  })

  it('validateRegisterForm 收集所有错误', () => {
    const errs = validateRegisterForm({
      email: 'bad',
      password: '123',
      nickname: '',
      work_years: '999',
    })
    expect(isFormValid(errs)).toBe(false)
    expect(errs.email).not.toBe('')
    expect(errs.password).not.toBe('')
    expect(errs.nickname).not.toBe('')
    expect(errs.work_years).not.toBe('')
  })
})
