<template>
  <AdminLayout>
    <div class="w-full flex flex-col lg:flex-row gap-6 select-none font-sans text-slate-800 min-h-[calc(100vh-120px)]">

      <!-- Left Pane / Intern Selection Sidebar (lg:w-72) -->
      <div class="w-full lg:w-72 bg-white/85 backdrop-blur-md rounded-xl border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col shrink-0 overflow-hidden">
        <div class="p-4 border-b border-[#F8BBD0] bg-[#FCE4EC]/40">
          <div class="relative">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-sm">search</span>
            <input 
              v-model="internSearchQuery"
              type="text" 
              placeholder="Cari Intern..."
              class="w-full pl-9 pr-3 py-2 bg-white rounded-lg border border-[#F8BBD0] text-xs font-semibold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 transition-all"
            />
          </div>
        </div>

        <div class="flex flex-col p-2 space-y-1.5 overflow-y-auto max-h-[500px] lg:max-h-none">
          <div 
            v-for="intern in filteredInterns" 
            :key="intern.id"
            @click="selectIntern(intern)"
            :class="[
              selectedIntern.id === intern.id 
                ? 'bg-[#F06292] text-white shadow-xs' 
                : 'bg-white hover:bg-[#FCE4EC]/40 text-[#1b1c1c]',
              'flex items-center gap-3 cursor-pointer px-3.5 py-3 rounded-lg transition-all border border-[#F8BBD0]'
            ]"
          >
            <img 
              :src="intern.avatar" 
              :alt="intern.name" 
              class="w-10 h-10 rounded-full object-cover border"
              :class="selectedIntern.id === intern.id ? 'border-white' : 'border-[#F8BBD0]'"
            />
            <div class="min-w-0 flex-1">
              <h4 class="font-bold text-xs truncate leading-tight">{{ intern.name }}</h4>
              <p class="text-[11px] truncate mt-0.5" :class="selectedIntern.id === intern.id ? 'text-white/80' : 'text-[#574146]'">
                {{ intern.dept }}
              </p>
            </div>
          </div>

          <div v-if="filteredInterns.length === 0" class="py-8 text-center text-xs text-[#8a7176] font-semibold">
            Intern tidak ditemukan.
          </div>
        </div>
      </div>

      <!-- Right Pane / Main Activity Content Area -->
      <div class="flex-1 flex flex-col gap-6 bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] overflow-hidden">
        
        <!-- Header & Intern Info -->
        <div class="flex flex-col gap-3 border-b border-[#F8BBD0] pb-6">
          <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <div class="flex items-center gap-3 flex-wrap">
              <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
                <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill" style="font-variation-settings: 'FILL' 1;">history_edu</span>
                Jurnal Aktivitas
              </h1>
              <span class="text-xs font-bold text-[#ab2c5d] bg-[#ffd9e4] px-3 py-1 rounded-full border border-[#F8BBD0]">
                {{ selectedDateFormatted }}
              </span>
            </div>

            <button 
              @click="openAddActivityModal"
              class="bg-[#F06292] hover:bg-[#ab2c5d] text-white font-bold text-xs px-4 py-2.5 rounded-lg transition-all flex items-center gap-2 cursor-pointer border-0 shadow-xs shrink-0"
            >
              <span class="material-symbols-outlined text-base">add</span>
              <span>Tambah Aktivitas</span>
            </button>
          </div>

          <div class="flex items-center gap-4 flex-wrap text-xs text-[#574146] pt-1">
            <span class="font-bold text-[#1b1c1c] text-sm">{{ selectedIntern.name }}</span>
            <span class="text-slate-300">•</span>
            <span class="font-medium text-[#574146]">{{ selectedIntern.dept }}</span>
            <span class="text-slate-300">•</span>
            <div class="flex items-center gap-1 text-[#1B5E20] bg-[#E8F5E9] px-2.5 py-0.5 rounded-full border border-[#A5D6A7] font-bold text-[11px]">
              <span class="material-symbols-outlined text-sm">how_to_reg</span>
              <span>HADIR (08:30 AM)</span>
            </div>
          </div>
        </div>

        <!-- Activity Table Section (Waktu Upload Otomatis Saat Diunggah) -->
        <div class="flex flex-col gap-4">
          <h3 class="font-bold text-sm text-[#1b1c1c] uppercase tracking-wider">Log Aktivitas Harian</h3>
          
          <div class="overflow-x-auto rounded-xl border border-[#F8BBD0] bg-white">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-[#FCE4EC] border-b border-[#F8BBD0] text-[11px] font-bold text-[#574146] uppercase tracking-wider">
                  <th class="py-3 px-4 w-36">Waktu Upload</th>
                  <th class="py-3 px-4">Deskripsi Aktivitas</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#F8BBD0] text-xs text-[#1b1c1c]">
                <tr 
                  v-for="act in selectedInternActivities" 
                  :key="act.id"
                  class="hover:bg-[#FCE4EC]/30 transition-colors"
                >
                  <td class="py-4 px-4 font-mono font-bold text-[#574146] whitespace-nowrap align-top flex items-center gap-1.5">
                    <span class="material-symbols-outlined text-xs text-[#f06292]">schedule</span>
                    <span>{{ act.time }}</span>
                  </td>
                  <td class="py-4 px-4">
                    <div class="flex items-start gap-3">
                      <div class="w-8 h-8 rounded-lg bg-[#ffd9e4] text-[#ab2c5d] flex items-center justify-center shrink-0 border border-[#F8BBD0]">
                        <span class="material-symbols-outlined text-base">{{ act.icon || 'task' }}</span>
                      </div>
                      <div>
                        <h4 class="font-bold text-xs text-[#1b1c1c] mb-0.5">{{ act.title }}</h4>
                        <p class="text-[11px] text-[#574146] leading-relaxed">{{ act.description }}</p>
                      </div>
                    </div>
                  </td>
                </tr>

                <tr v-if="selectedInternActivities.length === 0">
                  <td colspan="2" class="py-12 text-center text-[#8a7176]">
                    <span class="material-symbols-outlined text-4xl block mb-2 opacity-50">pending_actions</span>
                    <span class="text-xs font-semibold">Belum ada aktivitas yang diunggah hari ini.</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Documentation & Photos Section -->
        <div class="flex flex-col gap-4 pt-2">
          <div class="flex justify-between items-center">
            <h3 class="font-bold text-sm text-[#1b1c1c] uppercase tracking-wider">Dokumentasi &amp; Bukti Foto</h3>
            <button 
              @click="openAddPhotoModal"
              class="flex items-center gap-1.5 text-[#ab2c5d] font-bold text-xs hover:bg-[#FCE4EC] px-3.5 py-1.5 rounded-lg transition-colors border border-[#F8BBD0] cursor-pointer bg-white"
            >
              <span class="material-symbols-outlined text-sm">add_photo_alternate</span> 
              <span>Upload File Baru</span>
            </button>
          </div>

          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div v-for="(photo, index) in selectedInternPhotos" :key="index" class="flex flex-col gap-1.5">
              <div 
                @click="previewPhoto(photo.url)"
                class="w-full aspect-video rounded-lg overflow-hidden border border-[#F8BBD0] relative group cursor-pointer bg-slate-100"
              >
                <img :src="photo.url" :alt="photo.caption" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center text-white">
                  <span class="material-symbols-outlined text-xl">zoom_in</span>
                </div>
              </div>
              <div class="text-[10px] font-bold text-[#574146] text-center bg-[#FCE4EC]/50 p-1.5 rounded border border-[#F8BBD0] flex items-center justify-center gap-1">
                <span class="material-symbols-outlined text-[12px] text-[#f06292]">schedule</span>
                <span>{{ photo.caption }}</span>
              </div>
            </div>

            <div v-if="selectedInternPhotos.length === 0" class="col-span-2 md:col-span-4 py-8 border border-dashed border-[#F8BBD0] rounded-xl text-center text-xs text-[#8a7176]">
              Belum ada foto dokumentasi diunggah.
            </div>
          </div>
        </div>

      </div>

    </div>

    <!-- ===== MODAL TAMBAH AKTIVITAS (Waktu Terdeteksi Otomatis Saat Diunggah) ===== -->
    <transition name="fade">
      <div v-if="addModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" @click.self="addModalOpen = false">
        <div class="w-full max-w-lg bg-white rounded-xl border border-[#F8BBD0] shadow-2xl p-6 flex flex-col gap-4 max-h-[90vh] overflow-y-auto">
          <div class="flex justify-between items-center border-b border-[#F8BBD0] pb-3">
            <h3 class="font-bold text-[#1b1c1c] text-base flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">add_task</span>
              Tambah Aktivitas Intern
            </h3>
            <button @click="addModalOpen = false" class="p-1 rounded-full hover:bg-[#FCE4EC] text-[#574146] border-0 bg-transparent cursor-pointer">
              <span class="material-symbols-outlined text-lg">close</span>
            </button>
          </div>

          <form @submit.prevent="saveNewActivity" class="flex flex-col gap-4">
            <!-- Automated Upload Time Info Badge -->
            <div class="flex items-center justify-between p-3 rounded-lg bg-[#FCE4EC] border border-[#F8BBD0]">
              <div class="flex items-center gap-2 text-xs font-bold text-[#ab2c5d]">
                <span class="material-symbols-outlined text-base">schedule</span>
                <span>Waktu Upload (Otomatis):</span>
              </div>
              <span class="text-xs font-mono font-bold text-[#1b1c1c] bg-white px-2.5 py-1 rounded border border-[#F8BBD0]">
                {{ liveUploadTime }}
              </span>
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-bold text-[#574146] uppercase">Judul Aktivitas <span class="text-[#f06292]">*</span></label>
              <input type="text" v-model="newAct.title" required placeholder="Contoh: Research Design System" class="w-full px-3.5 py-2 border border-[#F8BBD0] rounded-lg text-xs font-semibold focus:outline-none focus:border-[#f06292]" />
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-bold text-[#574146] uppercase">Deskripsi Detail</label>
              <textarea v-model="newAct.description" rows="3" placeholder="Jelaskan detail tugas/kegiatan yang diunggah..." class="w-full px-3.5 py-2 border border-[#F8BBD0] rounded-lg text-xs font-medium resize-none focus:outline-none focus:border-[#f06292]"></textarea>
            </div>

            <div class="flex gap-3 justify-end pt-3 border-t border-[#F8BBD0]">
              <button type="button" @click="addModalOpen = false" class="py-2 px-4 border border-[#F8BBD0] hover:bg-[#FCE4EC] text-xs font-bold text-[#574146] rounded-lg cursor-pointer bg-white">Batal</button>
              <button type="submit" class="py-2 px-5 bg-[#ab2c5d] hover:bg-[#8b0e45] text-white text-xs font-bold rounded-lg cursor-pointer border-0 shadow-xs">Unggah Aktivitas</button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- ===== MODAL UPLOAD FILE BARU (Waktu Deteksi Otomatis Saat Upload) ===== -->
    <transition name="fade">
      <div v-if="addPhotoModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" @click.self="addPhotoModalOpen = false">
        <div class="w-full max-w-md bg-white rounded-xl border border-[#F8BBD0] shadow-2xl p-6 flex flex-col gap-4">
          <div class="flex justify-between items-center border-b border-[#F8BBD0] pb-3">
            <h3 class="font-bold text-[#1b1c1c] text-base flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">add_photo_alternate</span>
              Upload File Baru
            </h3>
            <button @click="addPhotoModalOpen = false" class="p-1 rounded-full hover:bg-[#FCE4EC] text-[#574146] border-0 bg-transparent cursor-pointer">
              <span class="material-symbols-outlined text-lg">close</span>
            </button>
          </div>

          <form @submit.prevent="saveNewPhoto" class="flex flex-col gap-4">
            <!-- Automated Upload Time Info Badge -->
            <div class="flex items-center justify-between p-3 rounded-lg bg-[#FCE4EC] border border-[#F8BBD0]">
              <div class="flex items-center gap-2 text-xs font-bold text-[#ab2c5d]">
                <span class="material-symbols-outlined text-base">schedule</span>
                <span>Waktu Unggah File:</span>
              </div>
              <span class="text-xs font-mono font-bold text-[#1b1c1c] bg-white px-2.5 py-1 rounded border border-[#F8BBD0]">
                {{ liveFullUploadTime }}
              </span>
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-[10px] font-bold text-[#574146] uppercase">Pilih Foto / Dokumen</label>
              <div class="border-2 border-dashed border-[#F8BBD0] rounded-xl p-6 text-center bg-[#FCE4EC]/20 hover:bg-[#FCE4EC]/40 transition-colors cursor-pointer relative">
                <input type="file" @change="handleFileSelected" accept="image/*" class="absolute inset-0 opacity-0 cursor-pointer w-full h-full" />
                <span class="material-symbols-outlined text-3xl text-[#f06292] block mb-1">cloud_upload</span>
                <span class="text-xs font-bold text-[#ab2c5d] block">Klik untuk memilih foto</span>
                <span class="text-[10px] text-[#8a7176] block mt-0.5">Format JPG, PNG, WEBP (Maks 5MB)</span>
              </div>
            </div>

            <div class="flex gap-3 justify-end pt-3 border-t border-[#F8BBD0]">
              <button type="button" @click="addPhotoModalOpen = false" class="py-2 px-4 border border-[#F8BBD0] hover:bg-[#FCE4EC] text-xs font-bold text-[#574146] rounded-lg cursor-pointer bg-white">Batal</button>
              <button type="submit" class="py-2 px-5 bg-[#ab2c5d] hover:bg-[#8b0e45] text-white text-xs font-bold rounded-lg cursor-pointer border-0 shadow-xs">Upload Foto</button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- ===== PHOTO PREVIEW MODAL ===== -->
    <transition name="fade">
      <div v-if="previewPhotoUrl" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/80 backdrop-blur-md" @click.self="previewPhotoUrl = ''">
        <div class="relative max-w-3xl w-full bg-white rounded-xl overflow-hidden shadow-2xl p-2">
          <img :src="previewPhotoUrl" alt="Preview Foto" class="w-full h-auto max-h-[80vh] object-contain rounded-lg" />
          <button @click="previewPhotoUrl = ''" class="absolute top-4 right-4 bg-black/60 text-white rounded-full p-1.5 hover:bg-black border-0 cursor-pointer">
            <span class="material-symbols-outlined text-lg">close</span>
          </button>
        </div>
      </div>
    </transition>

  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import AdminLayout from '@/layouts/AdminLayout.vue'

const selectedDateFormatted = ref('24 Okt 2026')
const internSearchQuery = ref('')
const addModalOpen = ref(false)
const addPhotoModalOpen = ref(false)
const previewPhotoUrl = ref('')
const tempPhotoUrl = ref('')

const getFormattedTimeNow = () => {
  const now = new Date()
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  return `${hours}:${minutes} WITA`
}

const getFormattedDateTimeNow = () => {
  const now = new Date()
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agt', 'Sep', 'Okt', 'Nov', 'Des']
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  return `${now.getDate()} ${months[now.getMonth()]} ${now.getFullYear()} - ${hours}:${minutes} WITA`
}

const liveUploadTime = ref(getFormattedTimeNow())
const liveFullUploadTime = ref(getFormattedDateTimeNow())

const internsList = ref([
  {
    id: 1,
    name: 'Hikma',
    dept: 'Product Design',
    avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuCkAeHHM_k62vY_e9WMM24DUFSIsSWQ5bTTpPRut4bOo4NvLtKUGrcDsvaHF70IMj7S4gA_v3LhFG20gZHOa5SqxaCC4cbK75b9PizAMwXR9JD6otrQXUtbPcs6DVT1UTbdeXhryc-6rtWS0xC9fx4zC8S6UPMZkdqlpADyIzHZ8f1bw8ayWH2xE7ac8jDZLlt0YxSwN_Cwf114LG1qWu0i-qFWCRir6q4WIDH-RGpLM-yxBkQ3MzPX'
  },
  {
    id: 2,
    name: 'Budi Santoso',
    dept: 'Frontend Dev',
    avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuBBeE2BxTkb_pP78ApQKfBxJuF4rGCM8PFfor01cokcAv0UmO_EZRd1vo14wdFt6KZK1hVx1pUgTqznj2cJXT-hfiarQqF55IXOSPoNR6WrOCx8RJ2j_wDgoCux7nqsyGxt94WGoPpU-52q1M2cAmGEJRJpVoo2Pp3VtD8PwrOQqBid4ReCOYB7zbdYt1Qz7Q8c61u5qhEI819m1piVINhN-xODlec4weN9igndknR5ni0TUiuxkd4J'
  },
  {
    id: 3,
    name: 'Ayu Diah',
    dept: 'Backend Dev',
    avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuCrUbofspjLB4E8PIaO4bqEVb-Nv9dw4clsPFGnVgLRnDpjStNPrpX-vP1CGWxHkcwAH0E2okmLtW5PqyKheQEhRkz77w8-PSGa0J_da_2EOaC3ICqOHjxILKOwxU72vePWil6oxV7Ect1DLPJg0q7ZU9aGttmbBTuDoyRZuoS6w3K3wisf8ZrZ0iSwEgeCd9fDM7ZHiof8VDX8Wq7xtbxUgKRnZsSINkXg79i8xBL4RpIsRkQfIAN0'
  }
])

const selectedIntern = ref(internsList.value[0])

const filteredInterns = computed(() => {
  if (!internSearchQuery.value) return internsList.value
  return internsList.value.filter(i => 
    i.name.toLowerCase().includes(internSearchQuery.value.toLowerCase()) ||
    i.dept.toLowerCase().includes(internSearchQuery.value.toLowerCase())
  )
})

const activitiesMap = ref<Record<number, any[]>>({
  1: [
    { id: 101, time: '11:30 WITA', title: 'Research Design System', description: 'Mempelajari dan menyusun pedoman untuk Design System LOPI-Q yang baru.', icon: 'design_services' },
    { id: 102, time: '14:30 WITA', title: 'Meeting Coordination', description: 'Sinkronisasi progres mingguan dengan tim developer terkait implementasi UI.', icon: 'groups' },
    { id: 103, time: '16:15 WITA', title: 'Prototyping Jurnal Intern', description: 'Membuat high-fidelity prototype untuk halaman log aktivitas & presensi.', icon: 'pending_actions' }
  ],
  2: [
    { id: 201, time: '10:45 WITA', title: 'Slicing UI Vue 3', description: 'Menerapkan komponen Tailwind CSS & AdminLayout pada tampilan baru.', icon: 'code' }
  ],
  3: [
    { id: 301, time: '09:15 WITA', title: 'gRPC Endpoint Refactoring', description: 'Optimalisasi mikroservis activity-service & reporting-service.', icon: 'dns' }
  ]
})

const photosMap = ref<Record<number, any[]>>({
  1: [
    { url: 'https://lh3.googleusercontent.com/aida-public/AB6AXuDjG60y_Co42fGoabvcXy2ugA2BJ-BSKF3ZJu9k704bDXK1N4OX2YqDsOq8e7qQ68P6ICXDJ1kNHDX00gFrsmK3fwxQYzyv4eBHyam89DD4SavDZO2YrJyQFWBWx6ApLwcCHAUG13IQRQdDW5Xe0LqcVVppU_xOJgA0wzczjvJcBcfPDkQ4VjWe2Tj3nh-kNuNvqfxwfX-2icocII3EB2dV1c0GbJg40kGTualtnyWnS55v9VlbOUuX', caption: '24 Okt 2026 - 11:30 WITA' },
    { url: 'https://lh3.googleusercontent.com/aida-public/AB6AXuB2wcXtyQwXG8NimOcG925_9hwHR_t4sqUAIjA_F-RpfN_Suippiv7F3xptV5rqtbrZA3rR8Grjur9-lCqbL5Sti5hMdFRHGvl2S9-apvjW9MXwPz9ANKpQNWHlK7eSBD3QKZ2mAObgGZANyxfmmALhEU_E40kZLtXA1AZp5lqOEb5tSerlbnAgQ98JOdNfE8kcZpeSRggFXKgqu4cGuufVhhWb1i6Bq0tkYTaasIqQ-7jE4pKNSGTT', caption: '24 Okt 2026 - 14:30 WITA' }
  ],
  2: [],
  3: []
})

const selectedInternActivities = computed(() => {
  return activitiesMap.value[selectedIntern.value.id] || []
})

const selectedInternPhotos = computed(() => {
  return photosMap.value[selectedIntern.value.id] || []
})

const selectIntern = (intern: any) => {
  selectedIntern.value = intern
}

const previewPhoto = (url: string) => {
  previewPhotoUrl.value = url
}

const newAct = ref({
  title: '',
  description: ''
})

const openAddActivityModal = () => {
  liveUploadTime.value = getFormattedTimeNow()
  newAct.value = {
    title: '',
    description: ''
  }
  addModalOpen.value = true
}

const saveNewActivity = () => {
  const currentList = activitiesMap.value[selectedIntern.value.id] || []
  currentList.push({
    id: Date.now(),
    time: getFormattedTimeNow(),
    title: newAct.value.title,
    description: newAct.value.description,
    icon: 'task'
  })
  activitiesMap.value[selectedIntern.value.id] = currentList
  addModalOpen.value = false
}

const openAddPhotoModal = () => {
  liveFullUploadTime.value = getFormattedDateTimeNow()
  tempPhotoUrl.value = 'https://lh3.googleusercontent.com/aida-public/AB6AXuDjG60y_Co42fGoabvcXy2ugA2BJ-BSKF3ZJu9k704bDXK1N4OX2YqDsOq8e7qQ68P6ICXDJ1kNHDX00gFrsmK3fwxQYzyv4eBHyam89DD4SavDZO2YrJyQFWBWx6ApLwcCHAUG13IQRQdDW5Xe0LqcVVppU_xOJgA0wzczjvJcBcfPDkQ4VjWe2Tj3nh-kNuNvqfxwfX-2icocII3EB2dV1c0GbJg40kGTualtnyWnS55v9VlbOUuX'
  addPhotoModalOpen.value = true
}

const handleFileSelected = (event: any) => {
  const file = event.target.files[0]
  if (file) {
    tempPhotoUrl.value = URL.createObjectURL(file)
  }
}

const saveNewPhoto = () => {
  const currentPhotos = photosMap.value[selectedIntern.value.id] || []
  currentPhotos.push({
    url: tempPhotoUrl.value || 'https://lh3.googleusercontent.com/aida-public/AB6AXuDjG60y_Co42fGoabvcXy2ugA2BJ-BSKF3ZJu9k704bDXK1N4OX2YqDsOq8e7qQ68P6ICXDJ1kNHDX00gFrsmK3fwxQYzyv4eBHyam89DD4SavDZO2YrJyQFWBWx6ApLwcCHAUG13IQRQdDW5Xe0LqcVVppU_xOJgA0wzczjvJcBcfPDkQ4VjWe2Tj3nh-kNuNvqfxwfX-2icocII3EB2dV1c0GbJg40kGTualtnyWnS55v9VlbOUuX',
    caption: getFormattedDateTimeNow()
  })
  photosMap.value[selectedIntern.value.id] = currentPhotos
  addPhotoModalOpen.value = false
}
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
