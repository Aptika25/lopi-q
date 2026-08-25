<template>
  <div class="flex items-center justify-center gap-2">
    <input
      v-for="(_digit, index) in digits"
      :key="index"
      :ref="el => setInputRef(el, index)"
      v-model="digits[index]"
      type="text"
      maxlength="1"
      inputmode="numeric"
      pattern="[0-9]*"
      :disabled="disabled"
      class="w-10 h-12 text-center text-lg font-black bg-slate-50 border border-slate-300 focus:border-rose-600 focus:bg-white focus:ring-2 focus:ring-rose-500/20 rounded-xl text-slate-900 outline-none transition-all shadow-sm disabled:opacity-50 font-mono"
      @input="handleInput(index, $event)"
      @keydown="handleKeyDown(index, $event)"
      @paste="handlePaste"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    disabled?: boolean
  }>(),
  {
    modelValue: '',
    disabled: false
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'complete', value: string): void
}>()

const digits = ref<string[]>(['', '', '', '', '', ''])
const inputRefs = ref<HTMLInputElement[]>([])

function setInputRef(el: any, index: number) {
  if (el) {
    inputRefs.value[index] = el as HTMLInputElement
  }
}

watch(() => props.modelValue, (val) => {
  if (!val) {
    digits.value = ['', '', '', '', '', '']
    return
  }
  const chars = val.split('').slice(0, 6)
  for (let i = 0; i < 6; i++) {
    digits.value[i] = chars[i] || ''
  }
})

const handleInput = (index: number, event: Event) => {
  const target = event.target as HTMLInputElement
  const val = target.value.replace(/[^0-9]/g, '')
  digits.value[index] = val

  const fullVal = digits.value.join('')
  emit('update:modelValue', fullVal)

  if (val && index < 5) {
    inputRefs.value[index + 1]?.focus()
  }

  if (fullVal.length === 6) {
    emit('complete', fullVal)
  }
}

const handleKeyDown = (index: number, event: KeyboardEvent) => {
  if (event.key === 'Backspace' && !digits.value[index] && index > 0) {
    inputRefs.value[index - 1]?.focus()
  }
}

const handlePaste = (event: ClipboardEvent) => {
  event.preventDefault()
  const pasteData = event.clipboardData?.getData('text') || ''
  const numeric = pasteData.replace(/[^0-9]/g, '').slice(0, 6)
  
  if (numeric) {
    for (let i = 0; i < 6; i++) {
      digits.value[i] = numeric[i] || ''
    }
    const fullVal = digits.value.join('')
    emit('update:modelValue', fullVal)
    if (fullVal.length === 6) {
      emit('complete', fullVal)
    }
  }
}

onMounted(() => {
  inputRefs.value[0]?.focus()
})
</script>
