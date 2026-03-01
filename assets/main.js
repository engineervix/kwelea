// kwelea — main.js

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

// ── Search ──
const overlay = document.getElementById('searchOverlay');
const headerSearchInput = document.getElementById('searchInput');
const modalInput = document.getElementById('searchModalInput');
const searchResultsEl = document.getElementById('searchResults');

// FlexSearch index — loaded once on first modal open.
const __basePath = (window.__kwelea && window.__kwelea.basePath) || '';
let searchIndex = null;
let indexLoading = false;

async function loadSearchIndex() {
  if (searchIndex !== null || indexLoading) return;
  indexLoading = true;
  try {
    const res = await fetch(__basePath + '/search-index.json');
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const entries = await res.json();

    searchIndex = new FlexSearch.Document({
      tokenize: 'forward',
      document: {
        id: 'id',
        index: ['title', 'body'],
        store: ['title', 'path', 'body'],
      },
    });
    entries.forEach(e => searchIndex.add(e));
  } catch (err) {
    console.warn('[kwelea] search index failed to load:', err);
    searchIndex = null;
  } finally {
    indexLoading = false;
  }
}

function openSearch() {
  if (!overlay) return;
  overlay.classList.add('open');
  if (modalInput) {
    modalInput.focus();
    loadSearchIndex(); // lazy-load on first open
  }
}

function closeSearch() {
  if (!overlay) return;
  overlay.classList.remove('open');
  if (modalInput) modalInput.value = '';
  if (searchResultsEl) searchResultsEl.innerHTML = '<div class="search-empty">Type to search…</div>';
  focusedIdx = -1;
}

if (headerSearchInput) headerSearchInput.addEventListener('click', openSearch);

if (overlay) {
  overlay.addEventListener('click', e => { if (e.target === overlay) closeSearch(); });
}

document.addEventListener('keydown', e => {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') { e.preventDefault(); openSearch(); }
  if (e.key === 'Escape' && overlay && overlay.classList.contains('open')) closeSearch();
});

// ── Search results rendering ──
function escapeHTML(s) {
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

// Extract a ~160-char snippet around the first query match, with terms highlighted.
function makeSnippet(body, query) {
  const queryWords = query.trim().toLowerCase().split(/\s+/).filter(Boolean);
  const lower = body.toLowerCase();

  // Find position of first matching term.
  let matchPos = 0;
  for (const word of queryWords) {
    const idx = lower.indexOf(word);
    if (idx >= 0) { matchPos = idx; break; }
  }

  const winStart = Math.max(0, matchPos - 40);
  const winEnd = Math.min(body.length, winStart + 160);
  const raw = body.slice(winStart, winEnd);
  let snippet = (winStart > 0 ? '…' : '') + escapeHTML(raw) + (winEnd < body.length ? '…' : '');

  // Highlight each query term. Escape body first, then mark — so <mark> isn't escaped.
  queryWords.forEach(word => {
    const safe = escapeHTML(word).replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    snippet = snippet.replace(new RegExp(`(${safe})`, 'gi'), '<mark>$1</mark>');
  });

  return snippet;
}

function renderResults(query) {
  if (!searchResultsEl) return;
  if (!query.trim()) {
    searchResultsEl.innerHTML = '<div class="search-empty">Type to search…</div>';
    focusedIdx = -1;
    return;
  }
  if (!searchIndex) {
    searchResultsEl.innerHTML = '<div class="search-empty">Loading index…</div>';
    return;
  }

  const raw = searchIndex.search(query, { enrich: true, limit: 8 });

  // Deduplicate across title/body field results.
  const seen = new Set();
  const docs = [];
  raw.forEach(fieldResult => {
    fieldResult.result.forEach(({ id, doc }) => {
      if (!seen.has(id)) { seen.add(id); docs.push(doc); }
    });
  });

  if (docs.length === 0) {
    searchResultsEl.innerHTML = `<div class="search-empty">No results for "${escapeHTML(query)}"</div>`;
    focusedIdx = -1;
    return;
  }

  searchResultsEl.innerHTML = docs.map(doc => `
    <a href="${escapeHTML(__basePath + doc.path)}" class="search-result" onclick="closeSearch()">
      <div class="search-result-title">${escapeHTML(doc.title)}</div>
      <div class="search-result-path">${escapeHTML(doc.path)}</div>
      <div class="search-result-snippet">${makeSnippet(doc.body, query)}</div>
    </a>
  `).join('');

  focusedIdx = -1;
}

// ── Keyboard navigation in results ──
let focusedIdx = -1;

function moveFocus(delta) {
  if (!searchResultsEl) return;
  const items = searchResultsEl.querySelectorAll('.search-result');
  if (items.length === 0) return;
  focusedIdx = Math.max(-1, Math.min(items.length - 1, focusedIdx + delta));
  items.forEach((el, i) => el.classList.toggle('focused', i === focusedIdx));
  if (focusedIdx >= 0) items[focusedIdx].scrollIntoView({ block: 'nearest' });
}

if (modalInput) {
  modalInput.addEventListener('input', () => {
    renderResults(modalInput.value);
  });

  modalInput.addEventListener('keydown', e => {
    if (e.key === 'ArrowDown') { e.preventDefault(); moveFocus(1); }
    else if (e.key === 'ArrowUp') { e.preventDefault(); moveFocus(-1); }
    else if (e.key === 'Enter') {
      if (focusedIdx >= 0 && searchResultsEl) {
        const items = searchResultsEl.querySelectorAll('.search-result');
        if (items[focusedIdx]) items[focusedIdx].click();
      }
    }
  });
}

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
