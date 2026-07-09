import { describe, it, expect } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import EmpathyButton from '~/components/empathy/EmpathyButton.vue'

describe('EmpathyButton', () => {
  it('渲染共情计数与默认文案', async () => {
    const wrapper = await mountSuspended(EmpathyButton, {
      props: { count: 12, empathized: false },
    })
    expect(wrapper.text()).toContain('我懂你')
    expect(wrapper.text()).toContain('12')
    expect(wrapper.attributes('aria-pressed')).toBe('false')
  })

  it('已共情时显示已懂你并 aria-pressed=true', async () => {
    const wrapper = await mountSuspended(EmpathyButton, {
      props: { count: 13, empathized: true },
    })
    expect(wrapper.text()).toContain('已懂你')
    expect(wrapper.attributes('aria-pressed')).toBe('true')
  })
})
