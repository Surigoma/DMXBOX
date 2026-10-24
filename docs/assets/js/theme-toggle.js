(function () {
  'use strict';

  var storageKey = 'dmxbox-docs-theme';
  var preference = null;
  var media = window.matchMedia('(prefers-color-scheme: dark)');

  try {
    preference = window.localStorage.getItem(storageKey);
  } catch (_) {
    // Theme switching also works when browser storage is unavailable.
  }
  if (preference !== 'light' && preference !== 'dark') preference = null;

  function applyTheme() {
    var theme = preference || (media.matches ? 'dark' : 'light');
    window.jtd.setTheme(theme);
    document.documentElement.style.colorScheme = theme;
    var button = document.getElementById('theme-toggle');
    if (button) button.setAttribute('aria-pressed', String(theme === 'dark'));
  }

  // Runs in the head so each page starts with the selected color scheme.
  applyTheme();

  media.addEventListener('change', function () {
    if (!preference) applyTheme();
  });

  document.addEventListener('DOMContentLoaded', function () {
    var button = document.getElementById('theme-toggle');
    if (!button) return;
    button.hidden = false;
    applyTheme();
    button.addEventListener('click', function () {
      preference = window.jtd.getTheme() === 'dark' ? 'light' : 'dark';
      try {
        window.localStorage.setItem(storageKey, preference);
      } catch (_) {
        // Keep the selection for this page even if it cannot be saved.
      }
      applyTheme();
    });
  });
})();
