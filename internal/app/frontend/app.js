// Theme Management
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
}

document.addEventListener('DOMContentLoaded', () => {
    initTheme();
    loadTags();
    setupEventListeners();
});

document.getElementById('themeToggle').addEventListener('click', () => {
    const current = document.documentElement.getAttribute('data-theme');
    const next = current === 'light' ? 'dark' : 'light';
    document.documentElement.setAttribute('data-theme', next);
    localStorage.setItem('theme', next);
});

// Load emotion tags
async function loadTags() {
    try {
        const response = await fetch('/api/v1/tags');
        const data = await response.json();
        
        const select = document.getElementById('emotionSelect');
        data.tags.forEach(tag => {
            const option = document.createElement('option');
            option.value = tag;
            option.textContent = tag.charAt(0).toUpperCase() + tag.slice(1);
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading tags:', error);
        showNotification('Failed to load emotion tags', 'error');
    }
}

// Setup event listeners
function setupEventListeners() {
    const form = document.getElementById('quoteForm');
    const emotionSelect = document.getElementById('emotionSelect');
    const customTag = document.getElementById('customTag');
    
    // Clear custom input when emotion is selected
    emotionSelect.addEventListener('change', () => {
        if (emotionSelect.value) {
            customTag.value = '';
        }
    });
    
    // Clear emotion select when custom input is typed
    customTag.addEventListener('input', () => {
        if (customTag.value) {
            emotionSelect.value = '';
        }
    });
    
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const emotion = emotionSelect.value;
        const custom = customTag.value.trim();
        
        if (!emotion && !custom) {
            showNotification('Please select an emotion or enter a custom tag', 'error');
            return;
        }
        
        const tag = emotion || custom;
        await generateQuote(tag);
    });
}

// Generate quote
async function generateQuote(tag) {
    try {
        const btn = document.querySelector('.btn-primary');
        btn.disabled = true;
        btn.textContent = 'Generating...';
        
        const response = await fetch('/api/v1/quote', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ tag })
        });
        
        if (!response.ok) {
            throw new Error('Failed to generate quote');
        }
        
        const quote = await response.json();
        displayQuote(quote);
        addToHistory(quote);
        showNotification('Quote generated successfully!', 'success');
        
    } catch (error) {
        console.error('Error generating quote:', error);
        showNotification('Failed to generate quote', 'error');
    } finally {
        const btn = document.querySelector('.btn-primary');
        btn.disabled = false;
        btn.textContent = 'Generate Quote';
    }
}

// Display quote
function displayQuote(quote) {
    const display = document.getElementById('quoteDisplay');
    const tag = quote.tag || 'Unknown';
    const text = quote.quote || quote.quote_text || quote.text || 'No quote available';
    const author = quote.author || 'Unknown';
    const source = quote.source || '';
    
    display.innerHTML = `
        <div class="quote-content">
            <span class="quote-tag">${tag}</span>
            <blockquote class="quote-text">"${text}"</blockquote>
            <div class="quote-author">— ${author}</div>
            ${source ? `<div class="quote-source">${source}</div>` : ''}
            <div class="quote-actions">
                <button class="btn-secondary" onclick="copyQuote('${text.replace(/'/g, "\\'")}', '${author.replace(/'/g, "\\'")}')">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                    </svg>
                    Copy
                </button>
                <button class="btn-secondary" onclick="shareQuote('${text.replace(/'/g, "\\'")}', '${author.replace(/'/g, "\\'")}')">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="18" cy="5" r="3"></circle>
                        <circle cx="6" cy="12" r="3"></circle>
                        <circle cx="18" cy="19" r="3"></circle>
                        <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line>
                        <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line>
                    </svg>
                    Share
                </button>
            </div>
        </div>
    `;
}

// Add to history
function addToHistory(quote) {
    const historyList = document.getElementById('historyList');
    const item = document.createElement('div');
    item.className = 'history-item';
    
    const tag = quote.tag || 'Unknown';
    const text = quote.quote || quote.quote_text || quote.text || 'No quote available';
    
    item.innerHTML = `
        <div class="tag">${tag}</div>
        <div class="preview">${text}</div>
    `;
    
    item.addEventListener('click', () => displayQuote(quote));
    
    historyList.insertBefore(item, historyList.firstChild);
    
    // Keep only last 10 items
    while (historyList.children.length > 10) {
        historyList.removeChild(historyList.lastChild);
    }
}

// Copy quote
function copyQuote(text, author) {
    const fullText = `"${text}" — ${author}`;
    navigator.clipboard.writeText(fullText).then(() => {
        showNotification('Quote copied to clipboard!', 'success');
    }).catch(() => {
        showNotification('Failed to copy quote', 'error');
    });
}

// Share quote
function shareQuote(text, author) {
    const fullText = `"${text}" — ${author}`;
    
    if (navigator.share) {
        navigator.share({
            title: 'QuoteBox',
            text: fullText,
        }).catch(() => {
            // User cancelled share
        });
    } else {
        copyQuote(text, author);
    }
}

// Show notification
function showNotification(message, type = 'success') {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = `notification ${type} show`;
    
    setTimeout(() => {
        notification.classList.remove('show');
    }, 3000);
}
