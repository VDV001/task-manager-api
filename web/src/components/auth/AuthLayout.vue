<script setup lang="ts">
import { CheckSquare } from 'lucide-vue-next'
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-primary relative overflow-hidden">
    <!-- Gradient background effects -->
    <div class="absolute inset-0">
      <div
        class="absolute inset-0 opacity-100"
        style="
          background: radial-gradient(100% 100% at 0% 0%, rgb(46, 46, 46) 0%, rgb(0, 0, 0) 100%);
          mask: radial-gradient(
            125% 100% at 0% 0%,
            rgb(0, 0, 0) 0%,
            rgba(0, 0, 0, 0.22) 88%,
            rgba(0, 0, 0, 0) 100%
          );
        "
      />
      <div
        class="absolute inset-0 opacity-15"
        style="
          background: linear-gradient(rgb(99, 102, 241) 0%, rgba(99, 102, 241, 0) 100%);
          mask: linear-gradient(
            90deg,
            rgba(0, 0, 0, 0) 0%,
            rgb(0, 0, 0) 20%,
            rgba(0, 0, 0, 0) 36%,
            rgb(0, 0, 0) 55%,
            rgba(0, 0, 0, 0) 97%
          );
          transform: skewX(45deg);
        "
      />
    </div>
    <!-- Dot pattern -->
    <div
      class="absolute inset-0 opacity-20"
      style="
        background-image: radial-gradient(
          circle at 1px 1px,
          rgba(255, 255, 255, 0.3) 1px,
          transparent 0
        );
        background-size: 24px 24px;
      "
    />
    <!-- Content -->
    <div class="relative z-10 w-full max-w-sm">
      <div class="flex items-center justify-center gap-2 mb-8">
        <CheckSquare class="w-8 h-8 text-accent" />
        <span class="text-2xl font-bold text-text-primary">TaskFlow</span>
      </div>
      <div
        class="rounded-2xl bg-gradient-to-br from-white/[0.06] to-transparent backdrop-blur-sm border border-border p-8"
      >
        <slot />
      </div>
    </div>
  </div>
</template>
