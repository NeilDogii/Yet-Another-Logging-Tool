<script>
    import { onMount } from 'svelte';

    let logs = [];
    let error = null;
    let expandedRows = new Set();
    let paginationData = {
        total_logs: 0,
        total_pages: 0,
        logs_per_page: 0
    }

    let page = 0;
    
    async function fetchLogs() {
        try {
            const response = await fetch(`http://localhost:3000/api/logs?page=${page}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const result = await response.json();
            logs = result.data;
            paginationData = result.meta;
            expandedRows = new Set(); 
        } catch (err) {
            error = err.message;
        }
    }
    
    onMount(fetchLogs);

    function nextPage() {
        if (page < paginationData.total_pages - 1) {
            page++;
            fetchLogs();
        }
    }

    function prevPage() {
        if (page > 0) {
            page--;
            fetchLogs();
        }
    }

    function goToPage(pageNum) {
        page = pageNum;
        fetchLogs();
    }

    function toggleRow(id) {
        if (expandedRows.has(id)) {
            expandedRows.delete(id);
        } else {
            expandedRows.add(id);
        }
        expandedRows = expandedRows;
    }

    function getLevelColor(level) {
        const colors = {
            'ERROR': '#ef4444',
            'WARN': '#f59e0b',
            'INFO': '#3b82f6',
            'DEBUG': '#8b5cf6',
            'TRACE': '#6b7280'
        };
        return colors[level?.toUpperCase()] || '#6b7280';
    }
</script>

<div class="container">
    {#if error}
        <p class="error">Error: {error}</p>
    {:else if logs.length > 0}
        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Timestamp (MM/DD/YYYY)</th>
                    <th>Level</th>
                    <th>Message</th>
                    <th>Source</th>
                    <th>Environment</th>
                    <th>Hostname</th>
                    <th>Metadata</th>
                </tr>
            </thead>
            <tbody>
                {#each logs as log}
                    <tr>
                        <td>{String(log.ID)}</td>
                        <td>{new Date(log.Timestamp).toLocaleString()}</td>
                        <td>
                            <span class="level-badge" style="background-color: {getLevelColor(log.Level)}">
                                {log.Level}
                            </span>
                        </td>
                        <td>{log.Message || "No Message"}</td>
                        <td>{log.Source || "unknown"}</td>
                        <td>{log.Environment || "unknown"}</td>
                        <td>{log.Hostname || "unknown"}</td>
                        <td>
                            <button class="metadata-toggle" on:click={() => toggleRow(log.ID)}>
                                {expandedRows.has(log.ID) ? '▼' : '▶'} View
                            </button>
                        </td>
                    </tr>
                    {#if expandedRows.has(log.ID)}
                        <tr class="metadata-row">
                            <td colspan="8" style="text-align: left;">
                                <pre class="metadata-content">{JSON.stringify(log.Metadata, null, 2)}</pre>
                            </td>
                        </tr>
                    {/if}
                {/each}
            </tbody>
        </table>
        <div class="pagination">
            <div class="pagination-info">
                Showing page {page + 1} of {paginationData.total_pages} 
                ({paginationData.total_logs} total logs)
            </div>
            <div class="pagination-controls">
                <button 
                    class="pagination-btn" 
                    on:click={prevPage} 
                    disabled={page === 0}
                >
                    ← Previous
                </button>
                
                <div class="page-numbers">
                    {#each Array(paginationData.total_pages) as _, i}
                        {#if i === 0 || i === paginationData.total_pages - 1 || (i >= page - 2 && i <= page + 2)}
                            <button 
                                class="page-btn {i === page ? 'active' : ''}" 
                                on:click={() => goToPage(i)}
                            >
                                {i + 1}
                            </button>
                        {:else if i === page - 3 || i === page + 3}
                            <span class="ellipsis">...</span>
                        {/if}
                    {/each}
                </div>

                <button 
                    class="pagination-btn" 
                    on:click={nextPage} 
                    disabled={page >= paginationData.total_pages - 1}
                >
                    Next →
                </button>
            </div>
        </div>
    {:else}
        <p class="loading">Loading logs...</p>
    {/if}
</div>

<style>
    .container {
        width: 100dvw;
        height: 100vh;
        overflow: auto;
        background: #0f172a;
        padding: 20px;
        box-sizing: border-box;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        background: #1e293b;
        border-radius: 8px;
        overflow: hidden;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
    }

    thead {
        background: #334155;
        position: sticky;
        top: 0;
        z-index: 10;
    }

    th {
        padding: 16px;
        text-align: left;
        font-weight: 600;
        color: #f1f5f9;
        text-transform: uppercase;
        font-size: 12px;
        letter-spacing: 0.5px;
    }

    tbody tr {
        border-bottom: 1px solid #334155;
        transition: background-color 0.2s;
        text-align: center;
    }
    th {
        text-align: center;
    }
    tbody tr:hover {
        background: #2d3748;
    }

    td {
        padding: 16px;
        color: #cbd5e1;
        font-size: 14px;
    }

    .level-badge {
        display: inline-block;
        padding: 4px 12px;
        border-radius: 12px;

    }
    .error, .loading {
        color: #f1f5f9;
        text-align: center;
        padding: 20px;
    }

    .pagination {
        margin-top: 20px;
        display: flex;
        flex-direction: column;
        gap: 12px;
        align-items: center;
    }

    .pagination-info {
        color: #94a3b8;
        font-size: 14px;
    }

    .pagination-controls {
        display: flex;
        gap: 8px;
        align-items: center;
    }

    .pagination-btn {
        background: #3b82f6;
        color: white;
        border: none;
        padding: 8px 16px;
        border-radius: 6px;
        cursor: pointer;
        font-size: 14px;
        transition: background-color 0.2s;
        font-weight: 500;
    }

    .pagination-btn:hover:not(:disabled) {
        background: #2563eb;
    }

    .pagination-btn:disabled {
        background: #475569;
        cursor: not-allowed;
        opacity: 0.5;
    }

    .page-numbers {
        display: flex;
        gap: 4px;
    }

    .page-btn {
        background: #1e293b;
        color: #cbd5e1;
        border: 1px solid #334155;
        padding: 8px 12px;
        border-radius: 4px;
        cursor: pointer;
        font-size: 14px;
        transition: all 0.2s;
        min-width: 40px;
    }

    .page-btn:hover {
        background: #334155;
        border-color: #3b82f6;
    }

    .page-btn.active {
        background: #3b82f6;
        color: white;
        border-color: #3b82f6;
        font-weight: 600;
    }

    .ellipsis {
        color: #64748b;
        padding: 8px 4px;
        font-size: 14px;
    }

    .metadata-toggle {
        background: #3b82f6;
        color: white;
        border: none;
        padding: 6px 12px;
        border-radius: 4px;
        cursor: pointer;
        font-size: 12px;
        transition: background-color 0.2s;
    }

    .metadata-toggle:hover {
        background: #2563eb;
    }

    .metadata-row {
        background: #0f172a;
    }

    .metadata-content {
        margin: 0;
        padding: 16px;
        background: #1e293b;
        border-radius: 4px;
        color: #94a3b8;
        font-size: 12px;
        overflow-x: auto;
    }

    .error, .loading {
        color: #f1f5f9;
        text-align: center;
        padding: 20px;
    }
</style>
