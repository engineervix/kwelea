// kwelea — main.js
// Extracted from _dev/docs-ui.html and extended with copy-button injection.

// ── Theme toggle ──
const html = document.documentElement;
const themeBtn = document.getElementById('themeToggle');
const sunIcon = document.getElementById('sunIcon');
const moonIcon = document.getElementById('moonIcon');

const saved = localStorage.getItem('theme') ||
  (matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
setTheme(saved);

function setTheme(t) {
  html.setAttribute('data-theme', t);
  localStorage.setItem('theme', t);
  if (sunIcon) sunIcon.style.display = t === 'dark' ? 'block' : 'none';
  if (moonIcon) moonIcon.style.display = t === 'light' ? 'block' : 'none';
}

if (themeBtn) {
  themeBtn.addEventListener('click', () => {
    setTheme(html.getAttribute('data-theme') === 'dark' ? 'light' : 'dark');
  });
}

// ── Mobile sidebar ──
const sidebar = document.getElementById('sidebar');
const backdrop = document.getElementById('sidebarBackdrop');
const menuToggle = document.getElementById('menuToggle');

if (menuToggle && sidebar && backdrop) {
  menuToggle.addEventListener('click', () => {
    sidebar.classList.toggle('open');
    backdrop.classList.toggle('open');
  });

  backdrop.addEventListener('click', () => {
    sidebar.classList.remove('open');
    backdrop.classList.remove('open');
  });
}

// ── Search modal ──
const overlay = document.getElementById('searchOverlay');
const searchInput = document.getElementById('searchInput');

function openSearch() {
  if (!overlay) return;
  overlay.classList.add('open');
  const modalInput = document.getElementById('searchModalInput');
  if (modalInput) modalInput.focus();
}

function closeSearch() {
  if (overlay) overlay.classList.remove('open');
}

if (searchInput) searchInput.addEventListener('click', openSearch);

if (overlay) {
  overlay.addEventListener('click', e => {
    if (e.target === overlay) closeSearch();
  });
}

document.addEventListener('keydown', e => {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') { e.preventDefault(); openSearch(); }
  if (e.key === 'Escape') closeSearch();
});

// ── ToC scroll spy ──
const headings = document.querySelectorAll('.prose h2, .prose h3');
const tocLinks = document.querySelectorAll('.toc-link');

if (headings.length > 0 && tocLinks.length > 0) {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        tocLinks.forEach(l => l.classList.remove('active'));
        const link = document.querySelector(`.toc-link[href="#${entry.target.id}"]`);
        if (link) link.classList.add('active');
      }
    });
  }, { rootMargin: '-20% 0px -70% 0px' });

  headings.forEach(h => observer.observe(h));
}

// ── Code copy buttons ──
// Goldmark wraps highlighted code in .highlight > pre; plain code is just pre.
// We inject a copy button into every pre inside .prose.
document.querySelectorAll('.prose pre').forEach(pre => {
  const btn = document.createElement('button');
  btn.className = 'code-copy';
  btn.textContent = 'copy';
  btn.addEventListener('click', () => {
    const code = pre.querySelector('code');
    const text = code ? code.textContent : pre.textContent;
    navigator.clipboard.writeText(text).then(() => {
      btn.textContent = 'copied';
      setTimeout(() => { btn.textContent = 'copy'; }, 2000);
    });
  });
  pre.style.position = 'relative';
  pre.appendChild(btn);
});
