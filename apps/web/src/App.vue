<template>
  <div class="min-h-screen flex flex-col bg-background text-on-background selection:bg-primary-container selection:text-on-primary font-sans antialiased">
    <!-- Navbar Header (Hidden if route has hideNavbarFooter) -->
    <Navbar v-if="!route.meta.hideNavbarFooter" />

    <!-- Main View Outlet -->
    <main class="flex-1 w-full">
      <router-view />
    </main>

    <!-- Footer (Single Unified Footer) -->
    <footer 
      v-if="!route.meta.hideNavbarFooter"
      class="border-t border-outline-variant/30 bg-surface-container-lowest py-6 text-center text-xs text-on-surface-variant"
    >
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col sm:flex-row items-center justify-between gap-2">
        <div>&copy; 2026 <strong>LOPI-Q Kab. Bulukumba</strong>. Hak Cipta Dilindungi Undang-Undang.</div>
        <div class="text-[11px] text-on-surface-variant/70">
          Developed by 
          <a 
            href="https://diskominfo.bulukumbakab.go.id" 
            target="_blank" 
            rel="noopener noreferrer" 
            class="font-bold text-primary hover:underline transition-all"
          >Diskominfo dan Persandian Kab. Bulukumba</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from './stores/auth';
import Navbar from './components/Navbar.vue';

const route = useRoute();
const authStore = useAuthStore();

onMounted(() => {
  if (authStore.token && !authStore.user) {
    authStore.fetchProfile();
  }
});
</script>
