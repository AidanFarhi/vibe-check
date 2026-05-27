// ── Score color shift ────────────────────────────────
const SCORE_THEMES = [
  { min: 7.5, bg: '#0d1e18', border: '#1a3028', scoreColor: '#4ecfb0', labelColor: '#3a6a4a', streakColor: '#2a4a38' },
  { min: 4.5, bg: '#101828', border: '#1a2540', scoreColor: '#7aadff', labelColor: '#3a5a80', streakColor: '#2d4060' },
  { min: 2.5, bg: '#1a1018', border: '#301520', scoreColor: '#d98878', labelColor: '#6a3a3a', streakColor: '#4a2530' },
  { min: 0,   bg: '#1e0e10', border: '#3a1420', scoreColor: '#c44f4f', labelColor: '#6a2a2a', streakColor: '#4a1a20' },
];

function applyScoreTheme() {
  const card = document.getElementById('todayCard');
  if (!card) return;
  const score = parseFloat(document.getElementById('todayScore').textContent);
  const theme = SCORE_THEMES.find(t => score >= t.min);
  card.style.background = theme.bg;
  card.style.setProperty('--today-card-border', theme.border);
  document.getElementById('todayScore').style.color = theme.scoreColor;
  document.getElementById('todayLabel').style.color = theme.labelColor;
  document.getElementById('todayTitle').style.color = theme.streakColor;
  document.querySelector('.today-streak').style.color = theme.streakColor;
}

applyScoreTheme();
document.addEventListener('htmx:afterSettle', applyScoreTheme);

// ── Chart filters ────────────────────────────────────
function setTime(btn) {
  document.querySelectorAll('.time-btn').forEach(b => b.classList.remove('active'));
  btn.classList.add('active');
}

function setMetricPill(pill) {
  document.querySelectorAll('.pill').forEach(p => p.classList.remove('active'));
  pill.classList.add('active');
  const id = pill.dataset.metric;
  document.querySelectorAll('.metric-card').forEach(c => {
    c.classList.toggle('active', c.dataset.metric === id);
  });
  applyMetric(id);
}

function selectMetric(id) {
  document.querySelectorAll('.pill').forEach(p => {
    p.classList.toggle('active', p.dataset.metric === id);
  });
  document.querySelectorAll('.metric-card').forEach(c => {
    c.classList.toggle('active', c.dataset.metric === id);
  });
  applyMetric(id);
}

function applyMetric(id) {
  document.querySelectorAll('.mline').forEach(g => {
    const match = g.dataset.metric === id;
    g.style.opacity = match ? '1' : '0';
    g.querySelector('.mline-area').style.display = match ? '' : 'none';
  });
  document.getElementById('todayCard')?.classList.toggle('active', id === 'score');
}

// ── Modal ────────────────────────────────────────────
function openModal() {
  document.getElementById('modalOverlay').classList.add('open');
  document.body.style.overflow = 'hidden';
}

function closeModal() {
  const overlay = document.getElementById('modalOverlay');
  overlay.classList.remove('open');
  overlay.addEventListener('transitionend', () => {
    document.body.style.overflow = '';
  }, { once: true });
}

function handleOverlayClick(e) {
  if (e.target === document.getElementById('modalOverlay')) closeModal();
}

function updateVal(key, val) {
  document.getElementById('val-' + key).textContent = val;
}

applyMetric('score');

// ── Chart dot fix ────────────────────────────────────
// preserveAspectRatio="none" stretches x and y independently, turning circles
// into ovals on wide screens. Compensate by scaleX on each dot.
function fixChartDots() {
  const svg = document.querySelector('.chart-svg');
  if (!svg) return;
  const w = svg.getBoundingClientRect().width;
  if (!w) return;
  const scaleX = 390 / w; // yScale/xScale = (130/100) / (w/300)
  svg.querySelectorAll('circle').forEach(c => {
    c.style.transform = `scaleX(${scaleX})`;
  });
}

let _dotFixTimer;
window.addEventListener('resize', () => {
  clearTimeout(_dotFixTimer);
  _dotFixTimer = setTimeout(fixChartDots, 60);
});
document.addEventListener('htmx:afterSettle', fixChartDots);
fixChartDots();
