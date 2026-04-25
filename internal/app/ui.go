package app

import (
	"fmt"
	"ikoyhn/podcast-sponsorblock/internal/config"
	"net/url"
)

func buildWebUI(baseURL string) string {
	token := config.AppConfig.Authentication.Token
	tokenJS := "null"
	if token != "" {
		tokenJS = fmt.Sprintf("%q", url.QueryEscape(token))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>CleanCast</title>
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    background: #0f172a; color: #e2e8f0; min-height: 100vh;
    display: flex; align-items: center; justify-content: center; padding: 2rem;
  }
  .card {
    background: #1e293b; border-radius: 1rem; padding: 2.5rem;
    width: 100%%; max-width: 540px; box-shadow: 0 25px 50px rgba(0,0,0,.4);
  }
  h1 { font-size: 1.75rem; font-weight: 700; color: #f8fafc; margin-bottom: .25rem; }
  .subtitle { color: #94a3b8; font-size: .9rem; margin-bottom: 2rem; }
  label { display: block; font-size: .85rem; font-weight: 500; color: #cbd5e1; margin-bottom: .5rem; }
  input[type=text] {
    width: 100%%; padding: .75rem 1rem; border-radius: .5rem;
    border: 1px solid #334155; background: #0f172a; color: #f1f5f9;
    font-size: 1rem; outline: none; transition: border-color .2s;
  }
  input[type=text]:focus { border-color: #6366f1; }
  .hint { font-size: .78rem; color: #64748b; margin-top: .4rem; }
  .type-row { display: flex; gap: .75rem; margin: 1rem 0; }
  .type-btn {
    flex: 1; padding: .5rem; border-radius: .4rem; border: 1px solid #334155;
    background: #0f172a; color: #94a3b8; cursor: pointer; font-size: .85rem;
    transition: all .15s;
  }
  .type-btn.active { background: #6366f1; border-color: #6366f1; color: #fff; }
  .result { margin-top: 1.5rem; display: none; }
  .result-label { font-size: .8rem; font-weight: 600; color: #94a3b8; text-transform: uppercase; letter-spacing: .05em; margin-bottom: .5rem; }
  .url-box {
    display: flex; gap: .5rem; align-items: center;
    background: #0f172a; border: 1px solid #334155; border-radius: .5rem;
    padding: .75rem 1rem;
  }
  .url-text { flex: 1; font-family: monospace; font-size: .82rem; color: #a5b4fc; word-break: break-all; }
  .copy-btn {
    flex-shrink: 0; padding: .4rem .75rem; border-radius: .35rem;
    border: none; background: #6366f1; color: #fff; cursor: pointer;
    font-size: .8rem; font-weight: 500; transition: background .15s;
  }
  .copy-btn:hover { background: #4f46e5; }
  .copy-btn.copied { background: #10b981; }
</style>
</head>
<body>
<div class="card">
  <h1>CleanCast</h1>
  <p class="subtitle">Generate your podcast RSS feed URL</p>

  <label for="feedId">YouTube Playlist or Channel ID</label>
  <input type="text" id="feedId" placeholder="PLxxx… or UCxxx…" autocomplete="off" spellcheck="false">
  <p class="hint">Paste a playlist ID (starts with PL) or a channel ID (starts with UC)</p>

  <div class="type-row">
    <button class="type-btn active" id="btn-auto" onclick="setType('auto')">Auto-detect</button>
    <button class="type-btn" id="btn-playlist" onclick="setType('playlist')">Playlist</button>
    <button class="type-btn" id="btn-channel" onclick="setType('channel')">Channel</button>
  </div>

  <div class="result" id="result">
    <div class="result-label">Your RSS URL</div>
    <div class="url-box">
      <span class="url-text" id="feedUrl"></span>
      <button class="copy-btn" id="copyBtn" onclick="copyUrl()">Copy</button>
    </div>
  </div>
</div>

<script>
  const BASE_URL = %q;
  const TOKEN = %s;
  let forceType = 'auto';

  function setType(t) {
    forceType = t;
    ['auto','playlist','channel'].forEach(function(x) {
      document.getElementById('btn-'+x).classList.toggle('active', x === t);
    });
    updateUrl();
  }

  function detectEndpoint(id) {
    if (forceType === 'playlist') return 'rss';
    if (forceType === 'channel')  return 'channel';
    return id.startsWith('PL') ? 'rss' : 'channel';
  }

  function updateUrl() {
    const id = document.getElementById('feedId').value.trim();
    const result = document.getElementById('result');
    if (!id) { result.style.display = 'none'; return; }
    const endpoint = detectEndpoint(id);
    let u = BASE_URL + '/' + endpoint + '/' + encodeURIComponent(id);
    if (TOKEN) u += '?token=' + TOKEN;
    document.getElementById('feedUrl').textContent = u;
    result.style.display = 'block';
  }

  function copyUrl() {
    const text = document.getElementById('feedUrl').textContent;
    navigator.clipboard.writeText(text).then(function() {
      const btn = document.getElementById('copyBtn');
      btn.textContent = 'Copied!';
      btn.classList.add('copied');
      setTimeout(function() { btn.textContent = 'Copy'; btn.classList.remove('copied'); }, 1500);
    });
  }

  document.getElementById('feedId').addEventListener('input', updateUrl);
</script>
</body>
</html>`, baseURL, tokenJS)
}
