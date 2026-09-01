<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">

      <!-- Header Section -->
      <div class="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill" style="font-variation-settings: 'FILL' 1;">history_edu</span>
            Jurnal Aktivitas
          </h1>
          <p class="text-sm text-[#574146] mt-1">Rekapitulasi log aktivitas harian peserta magang LOPI-Q.</p>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <span class="text-xs font-bold text-[#ab2c5d] bg-[#ffd9e4] px-4 py-2 rounded-lg border border-[#F8BBD0] flex items-center gap-1.5 shadow-xs">
            <span class="material-symbols-outlined text-base">calendar_today</span>
            <span>{{ selectedDateFormatted }}</span>
          </span>
        </div>
      </div>

      <!-- Filters & Search Bar (Glass Card) -->
      <div class="bg-white/85 backdrop-blur-md rounded-xl p-4 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col md:flex-row gap-4 items-center justify-between">
        <div class="flex flex-col md:flex-row gap-3 w-full md:w-auto">
          <!-- Department Filter -->
          <div class="relative w-full md:w-60">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-sm">domain</span>
            <select 
              v-model="filterDept" 
              class="w-full pl-10 pr-8 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 appearance-none transition-all cursor-pointer"
            >
              <option value="">Semua Departemen / Divisi</option>
              <option value="Product Design">Product Design</option>
              <option value="Frontend Dev">Frontend Dev</option>
              <option value="Backend Dev">Backend Dev</option>
            </select>
          </div>
        </div>

        <!-- Search Input -->
        <div class="relative w-full md:w-72">
          <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-sm">search</span>
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Cari nama intern atau aktivitas..." 
            class="w-full pl-10 pr-4 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-medium text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 transition-all"
          />
        </div>
      </div>

      <!-- Main Activity Table Card -->
      <div class="bg-white/90 backdrop-blur-md rounded-xl border border-[#F8BBD0] overflow-hidden shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
        <div class="overflow-x-auto w-full">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-[#FCE4EC] border-b border-[#F8BBD0] text-[11px] font-bold text-[#574146] uppercase tracking-wider">
                <th class="py-3.5 px-6 min-w-[200px]">Nama Intern</th>
                <th class="py-3.5 px-6 w-36">Waktu Upload</th>
                <th class="py-3.5 px-6">Deskripsi Aktivitas</th>
              </tr>
            </thead>
            <tbody class="text-xs text-[#1b1c1c] bg-white divide-y divide-[#F8BBD0]">
              <tr 
                v-for="act in filteredActivities" 
                :key="act.id"
                class="hover:bg-[#FCE4EC]/30 transition-colors"
              >
                <!-- Intern Name & Dept -->
                <td class="py-4 px-6 min-w-[200px] align-top">
                  <div class="font-bold text-[#1b1c1c] text-xs leading-tight">{{ act.internName }}</div>
                  <div class="text-[10px] text-[#574146] font-medium mt-0.5">{{ act.dept }}</div>
                </td>

                <!-- Time -->
                <td class="py-4 px-6 font-mono font-bold text-[#574146] whitespace-nowrap align-top">
                  <div class="flex items-center gap-1.5">
                    <span class="material-symbols-outlined text-xs text-[#f06292]">schedule</span>
                    <span>{{ act.time }}</span>
                  </div>
                </td>

                <!-- Activity Details -->
                <td class="py-4 px-6 align-top">
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

              <!-- Empty State -->
              <tr v-if="filteredActivities.length === 0">
                <td colspan="3" class="py-6 text-center">
                  <div class="inline-flex items-center justify-center gap-1.5 text-xs text-slate-500 font-medium">
                    <span class="material-symbols-outlined text-slate-400" style="font-size: 14px !important;">info</span>
                    <span>Belum ada log aktivitas yang diunggah.</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import AdminLayout from '@/layouts/AdminLayout.vue'

const selectedDateFormatted = ref('24 Okt 2026')
const filterDept = ref('')
const searchQuery = ref('')

const allActivities = ref([
  { id: 101, internName: 'Hikma', dept: 'Product Design', time: '11:30 WITA', title: 'Research Design System', description: 'Mempelajari dan menyusun pedoman untuk Design System LOPI-Q yang baru.', icon: 'design_services' },
  { id: 102, internName: 'Hikma', dept: 'Product Design', time: '14:30 WITA', title: 'Meeting Coordination', description: 'Sinkronisasi progres mingguan dengan tim developer terkait implementasi UI.', icon: 'groups' },
  { id: 103, internName: 'Hikma', dept: 'Product Design', time: '16:15 WITA', title: 'Prototyping Jurnal Intern', description: 'Membuat high-fidelity prototype untuk halaman log aktivitas & presensi.', icon: 'pending_actions' },
  { id: 201, internName: 'Budi Santoso', dept: 'Frontend Dev', time: '10:45 WITA', title: 'Slicing UI Vue 3', description: 'Menerapkan komponen Tailwind CSS & AdminLayout pada tampilan baru.', icon: 'code' },
  { id: 301, internName: 'Ayu Diah', dept: 'Backend Dev', time: '09:15 WITA', title: 'gRPC Endpoint Refactoring', description: 'Optimalisasi mikroservis activity-service & reporting-service.', icon: 'dns' }
])

const filteredActivities = computed(() => {
  return allActivities.value.filter(act => {
    const matchesDept = !filterDept.value || act.dept === filterDept.value
    const matchesSearch = !searchQuery.value || 
      act.internName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      act.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      act.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchesDept && matchesSearch
  })
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>






