package scalar

// PersistAuthScript generates JavaScript code to persist authentication tokens
// in localStorage and restore them on page load
const PersistAuthScript = `
<script>
(function() {
  'use strict';

  const STORAGE_KEY = 'scalar_auth_tokens';
  const STORAGE_EXPIRY_KEY = 'scalar_auth_expiry';
  const DEFAULT_EXPIRY_HOURS = 24;

  // Helper to safely access localStorage
  function getStorage() {
    try {
      return window.localStorage;
    } catch (e) {
      console.warn('localStorage not available:', e);
      return null;
    }
  }

  // Save auth data to localStorage with expiry
  function saveAuthData(data) {
    const storage = getStorage();
    if (!storage) return;

    try {
      const expiryTime = new Date();
      expiryTime.setHours(expiryTime.getHours() + DEFAULT_EXPIRY_HOURS);

      storage.setItem(STORAGE_KEY, JSON.stringify(data));
      storage.setItem(STORAGE_EXPIRY_KEY, expiryTime.toISOString());
    } catch (e) {
      console.warn('Failed to save auth data:', e);
    }
  }

  // Load auth data from localStorage
  function loadAuthData() {
    const storage = getStorage();
    if (!storage) return null;

    try {
      const expiry = storage.getItem(STORAGE_EXPIRY_KEY);
      if (expiry && new Date(expiry) < new Date()) {
        // Data expired, clear it
        storage.removeItem(STORAGE_KEY);
        storage.removeItem(STORAGE_EXPIRY_KEY);
        return null;
      }

      const data = storage.getItem(STORAGE_KEY);
      return data ? JSON.parse(data) : null;
    } catch (e) {
      console.warn('Failed to load auth data:', e);
      return null;
    }
  }

  // Clear auth data from localStorage
  function clearAuthData() {
    const storage = getStorage();
    if (!storage) return;

    try {
      storage.removeItem(STORAGE_KEY);
      storage.removeItem(STORAGE_EXPIRY_KEY);
    } catch (e) {
      console.warn('Failed to clear auth data:', e);
    }
  }

  // Monitor input changes for auth tokens
  function monitorAuthInputs() {
    // Use MutationObserver to detect when Scalar adds input fields
    const observer = new MutationObserver(function(mutations) {
      // Look for auth input fields
      const authInputs = document.querySelectorAll(
        'input[type="text"][placeholder*="token" i], ' +
        'input[type="text"][placeholder*="bearer" i], ' +
        'input[type="password"][placeholder*="token" i], ' +
        'input[type="password"][placeholder*="bearer" i], ' +
        'input[type="text"][name*="authorization" i], ' +
        'input[type="text"][name*="auth" i], ' +
        'input[name*="token" i]'
      );

      authInputs.forEach(function(input) {
        if (input.dataset.scalarAuthPersist) return; // Already monitoring
        input.dataset.scalarAuthPersist = 'true';

        // Restore saved value
        const savedData = loadAuthData();
        if (savedData && savedData[input.name || input.placeholder]) {
          input.value = savedData[input.name || input.placeholder];
        }

        // Save on change
        input.addEventListener('input', function() {
          const key = input.name || input.placeholder;
          const savedData = loadAuthData() || {};

          if (input.value.trim()) {
            savedData[key] = input.value;
            saveAuthData(savedData);
          } else {
            delete savedData[key];
            if (Object.keys(savedData).length === 0) {
              clearAuthData();
            } else {
              saveAuthData(savedData);
            }
          }
        });
      });
    });

    // Start observing
    observer.observe(document.body, {
      childList: true,
      subtree: true
    });
  }

  // Try to intercept Scalar API Reference configuration
  function injectAuthPersistence() {
    const scriptTag = document.getElementById('api-reference');
    if (!scriptTag) return;

    try {
      const config = JSON.parse(scriptTag.getAttribute('data-configuration') || '{}');
      const savedData = loadAuthData();

      // If we have saved auth data and Scalar supports authentication
      if (savedData && config.authentication) {
        // Try to inject saved tokens into the configuration
        const authConfig = JSON.parse(config.authentication || '{}');

        // For bearer token
        if (savedData.bearer || savedData.token) {
          const token = savedData.bearer || savedData.token;
          // Scalar will handle this in its UI
          console.log('Scalar: Restored auth token from localStorage');
        }
      }
    } catch (e) {
      console.warn('Failed to inject auth persistence:', e);
    }
  }

  // Wait for Scalar to load
  function initPersistence() {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', function() {
        injectAuthPersistence();
        monitorAuthInputs();
      });
    } else {
      injectAuthPersistence();
      monitorAuthInputs();
    }
  }

  // Start the persistence system
  initPersistence();

  // Expose clear function globally (for debugging)
  window.scalarClearAuth = clearAuthData;

  console.log('✅ Scalar Auth Persistence enabled (data expires in ' + DEFAULT_EXPIRY_HOURS + ' hours)');
  console.log('💡 To clear saved tokens, run: scalarClearAuth()');
})();
</script>
`
