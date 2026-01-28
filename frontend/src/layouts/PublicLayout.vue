<!-- filepath: /home/kali/Desktop/social-network/frontend/src/layouts/PublicLayout.vue -->
<template>
  <div class="public-layout">
    <nav class="navbar">
      <div class="navbar-brand">
        <router-link to="/login" class="brutalist-logo">
          <div class="logo-grid">
            <div class="logo-square logo-c">C</div>
            <div class="logo-square logo-ampersand">&</div>
            <div class="logo-square logo-d">D</div>
          </div>
        </router-link>
      </div>
    </nav>
    <main class="main-content">
      <slot></slot>
    </main>
  </div>
</template>

<script>
export default {
  name: 'PublicLayout'
}
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
  --brutalist-black: #000;
  --brutalist-white: #fff;
  --brutalist-blue: #00a4ef;
}

.public-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--lavender-mist);
}

.navbar {
  background: white;
  border-bottom: 1px solid rgba(13, 19, 33, 0.08);
  padding: 16px 24px;
  box-shadow: 0 2px 8px rgba(13, 19, 33, 0.04);
}

.navbar-brand {
  display: flex;
  align-items: center;
  justify-content: center;
}

/* =========================== BRUTALIST C&D LOGO =========================== */
.brutalist-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  width: 120px;
  height: 60px;
  background-color: var(--brutalist-black);
  color: var(--brutalist-white);
  text-decoration: none;
  font-family: Arial, sans-serif;
  font-weight: bold;
  border: 3px solid var(--brutalist-white);
  outline: 3px solid var(--brutalist-black);
  box-shadow: 6px 6px 0 var(--brutalist-blue);
  transition: all 0.1s ease-out;
  padding: 0;
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
}

/* Shine effect on hover */
.brutalist-logo::before {
  content: "";
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.8),
    transparent
  );
  z-index: 1;
  transition: none;
  opacity: 0;
}

@keyframes slide {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}

.brutalist-logo:hover::before {
  opacity: 1;
  animation: slide 2s infinite;
}

.brutalist-logo:hover {
  transform: translate(-4px, -4px);
  box-shadow: 10px 10px 0 var(--brutalist-black);
  background-color: var(--brutalist-black);
  color: var(--brutalist-white);
}

.brutalist-logo:active {
  transform: translate(4px, 4px);
  box-shadow: 0px 0px 0 var(--brutalist-blue);
  background-color: var(--brutalist-white);
  color: var(--brutalist-black);
  border-color: var(--brutalist-black);
}

/* Logo Grid - 3 squares for C & D */
.logo-grid {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  position: relative;
  z-index: 2;
  transition: transform 0.2s ease-out;
}

.brutalist-logo:hover .logo-grid {
  transform: rotate(-5deg) scale(1.05);
}

.brutalist-logo:active .logo-grid {
  transform: rotate(5deg) scale(0.95);
}

.logo-square {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 900;
  color: var(--brutalist-white);
  text-transform: uppercase;
  letter-spacing: -0.05em;
}

.logo-c {
  background-color: #f25022; /* Microsoft red */
}

.logo-ampersand {
  background-color: #7fba00; /* Microsoft green */
  font-size: 16px;
}

.logo-d {
  background-color: #00a4ef; /* Microsoft blue */
}

/* Active state - invert colors */
.brutalist-logo:active .logo-square {
  color: var(--brutalist-black);
}

.main-content {
  flex: 1;
  padding: 32px 16px;
}

/* =========================== RESPONSIVE =========================== */
@media (max-width: 768px) {
  .navbar {
    padding: 12px 16px;
  }

  .brutalist-logo {
    width: 90px;
    height: 44px;
    border: 2px solid var(--brutalist-white);
    outline: 2px solid var(--brutalist-black);
    box-shadow: 4px 4px 0 var(--brutalist-blue);
  }

  .brutalist-logo:hover {
    box-shadow: 6px 6px 0 var(--brutalist-black);
  }

  .logo-square {
    width: 24px;
    height: 24px;
    font-size: 16px;
  }

  .logo-ampersand {
    font-size: 12px;
  }
}

@media (max-width: 480px) {
  .brutalist-logo {
    width: 80px;
    height: 40px;
  }

  .logo-square {
    width: 20px;
    height: 20px;
    font-size: 14px;
  }

  .logo-ampersand {
    font-size: 11px;
  }

  .main-content {
    padding: 24px 12px;
  }
}
</style>