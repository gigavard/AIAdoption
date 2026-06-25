const API_URL = '/todos';

async function fetchTodos() {
    try {
        const response = await fetch(API_URL);
        if (!response.ok) throw new Error('Failed to fetch todos');
        return await response.json() || [];
    } catch (error) {
        showError('Failed to load todos');
        return [];
    }
}

async function createTodo(title, content) {
    try {
        const response = await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title, content })
        });

        if (response.status === 201) {
            return await response.json();
        } else {
            const error = await response.json();
            throw new Error(error.detail || 'Failed to create todo');
        }
    } catch (error) {
        showError(error.message);
        throw error;
    }
}

async function updateTodo(id, updates) {
    try {
        const response = await fetch(`${API_URL}/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(updates)
        });

        if (!response.ok) throw new Error('Failed to update todo');
        return await response.json();
    } catch (error) {
        showError(error.message);
        throw error;
    }
}

async function completeTodo(id) {
    try {
        const response = await fetch(`${API_URL}/${id}/complete`, {
            method: 'PATCH'
        });

        if (!response.ok) throw new Error('Failed to complete todo');
        return await response.json();
    } catch (error) {
        showError(error.message);
        throw error;
    }
}

async function deleteTodo(id) {
    try {
        const response = await fetch(`${API_URL}/${id}`, {
            method: 'DELETE'
        });

        if (response.status !== 204) throw new Error('Failed to delete todo');
    } catch (error) {
        showError(error.message);
        throw error;
    }
}

function renderTodos(todos) {
    const list = document.getElementById('todosList');
    const empty = document.getElementById('emptyState');

    if (!todos || todos.length === 0) {
        list.innerHTML = '';
        empty.classList.remove('hidden');
        return;
    }

    empty.classList.add('hidden');
    list.innerHTML = todos.map(todo => `
        <div class="flex items-center gap-4 p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
            <button 
                onclick="toggleTodo(${todo.id})"
                class="flex-shrink-0 w-6 h-6 rounded border-2 ${todo.status === 'completed' ? 'bg-green-500 border-green-500' : 'border-gray-300'} transition-colors"
            >
                ${todo.status === 'completed' ? '✓' : ''}
            </button>
            <div class="flex-1 min-w-0">
                <p class="${todo.status === 'completed' ? 'line-through text-gray-400' : 'text-gray-800'} font-medium">
                    ${escapeHtml(todo.title)}
                </p>
                ${todo.content ? `<p class="text-gray-600 text-sm">${escapeHtml(todo.content)}</p>` : ''}
                <p class="text-gray-400 text-xs mt-1">${new Date(todo.created_at).toLocaleString()}</p>
            </div>
            <button 
                onclick="removeTodo(${todo.id})"
                class="flex-shrink-0 px-3 py-1 text-red-500 hover:bg-red-50 rounded transition-colors"
            >
                Delete
            </button>
        </div>
    `).join('');
}

async function toggleTodo(id) {
    const todos = await fetchTodos();
    const todo = todos.find(t => t.id === id);
    
    if (todo.status === 'completed') {
        await updateTodo(id, { status: 'pending' });
    } else {
        await completeTodo(id);
    }
    
    render();
}

async function removeTodo(id) {
    if (confirm('Are you sure?')) {
        await deleteTodo(id);
        render();
    }
}

function showError(message) {
    const msg = document.getElementById('errorMsg');
    msg.textContent = message;
    msg.classList.remove('hidden');
    setTimeout(() => msg.classList.add('hidden'), 5000);
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

async function handleAddTodo() {
    const title = document.getElementById('todoTitle').value.trim();
    const content = document.getElementById('todoContent').value.trim();

    if (!title) {
        showError('Please enter a title');
        return;
    }

    try {
        await createTodo(title, content);
        document.getElementById('todoTitle').value = '';
        document.getElementById('todoContent').value = '';
        render();
    } catch (error) {
        // Error already shown by createTodo
    }
}

async function render() {
    const todos = await fetchTodos();
    renderTodos(todos);
}

document.getElementById('addTodoBtn').addEventListener('click', handleAddTodo);
document.getElementById('todoTitle').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') handleAddTodo();
});

render();
