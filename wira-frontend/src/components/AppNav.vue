<template>
    <nav class="fixed bottom-0 w-full bg-white shadow-lg md:hidden">
      <div class="flex justify-around p-4">
        <router-link to="/dashboard">Dashboard</router-link>
        <router-link to="/profile">Profile</router-link>
      </div>
    </nav>
  </template>
  
  <style>
  /* Stack table rows on mobile */
  @media (max-width: 640px) {
    table td {
      @apply block w-full text-center;
    }
    
    table th {
      @apply hidden;
    }
  }
  </style>