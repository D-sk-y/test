import { useCallback, useEffect, useState } from 'react'

interface Task {
  id: number
  name: string
  status: string
  created_at: string
  updated_at: string
}

const API_BASE = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`${res.status} ${res.statusText}: ${text}`)
  }
  // 部分接口(如 PUT)无返回体
  if (res.status === 204) return undefined as T
  return res.json() as Promise<T>
}

function App() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  // 新建表单
  const [newName, setNewName] = useState('')
  const [newStatus, setNewStatus] = useState('pending')

  // v2 hello 表单
  const [helloMsg, setHelloMsg] = useState('hello worker')
  const [helloResult, setHelloResult] = useState('')

  const fetchTasks = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const data = await request<Task[]>('/v1/tasks')
      setTasks(data)
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchTasks()
  }, [fetchTasks])

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!newName.trim()) return
    try {
      const created = await request<Task>('/v1/tasks', {
        method: 'POST',
        body: JSON.stringify({ name: newName.trim(), status: newStatus }),
      })
      setTasks((prev) => [created, ...prev])
      setNewName('')
      setNewStatus('pending')
      setMessage(`已创建任务 #${created.id}`)
    } catch (e) {
      setError((e as Error).message)
    }
  }

  async function handleUpdate(task: Task) {
    try {
      await request<void>(`/v1/tasks/${task.id}`, {
        method: 'PUT',
        body: JSON.stringify({ id: task.id, name: task.name, status: task.status }),
      })
      setMessage(`已更新任务 #${task.id}`)
    } catch (e) {
      setError((e as Error).message)
    }
  }

  function toggleStatus(task: Task) {
    const next = task.status === 'done' ? 'pending' : 'done'
    const updated = { ...task, status: next }
    setTasks((prev) => prev.map((t) => (t.id === task.id ? updated : t)))
    handleUpdate(updated)
  }

  async function handleHello(e: React.FormEvent) {
    e.preventDefault()
    try {
      const data = await request<{ status: string }>('/v2/hello', {
        method: 'POST',
        body: JSON.stringify({ name: 'web', msg: helloMsg }),
      })
      setHelloResult(JSON.stringify(data))
      setMessage('v2 /hello 调用成功，消息已发到 Bus')
    } catch (e) {
      setError((e as Error).message)
    }
  }

  return (
    <div className="container">
      <header>
        <h1>Task Manager</h1>
        <p>React + TypeScript + Vite 前端，调用 Go 后端 RESTful v1 / RPC v2 接口</p>
      </header>

      <div className="status-bar">
        <span>后端: http://localhost:18080 (经 Vite 代理)</span>
        {loading && <span className="tag tag-loading">加载中…</span>}
        {message && <span className="tag tag-ok">{message}</span>}
        {error && <span className="tag tag-err">{error}</span>}
      </div>

      <section className="card">
        <h2>新建任务 (POST /api/v1/tasks)</h2>
        <form onSubmit={handleCreate} className="inline-form">
          <input
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            placeholder="任务名称"
            required
          />
          <select value={newStatus} onChange={(e) => setNewStatus(e.target.value)}>
            <option value="pending">pending</option>
            <option value="done">done</option>
          </select>
          <button type="submit">创建</button>
        </form>
      </section>

      <section className="card">
        <h2>
          任务列表 (GET /api/v1/tasks)
          <button className="ghost" onClick={fetchTasks} disabled={loading}>
            刷新
          </button>
        </h2>
        {tasks.length === 0 ? (
          <p className="empty">暂无任务，先创建一个吧。</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>名称</th>
                <th>状态</th>
                <th>创建时间</th>
                <th>更新时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {tasks.map((t) => (
                <tr key={t.id}>
                  <td>{t.id}</td>
                  <td>{t.name}</td>
                  <td>
                    <span className={`status ${t.status}`}>{t.status}</span>
                  </td>
                  <td>{t.created_at}</td>
                  <td>{t.updated_at}</td>
                  <td>
                    <button onClick={() => toggleStatus(t)}>
                      切换为 {t.status === 'done' ? 'pending' : 'done'}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

      <section className="card">
        <h2>RPC v2 消息 (POST /api/v2/hello)</h2>
        <form onSubmit={handleHello} className="inline-form">
          <input value={helloMsg} onChange={(e) => setHelloMsg(e.target.value)} placeholder="消息内容" />
          <button type="submit">发送到 Bus</button>
        </form>
        {helloResult && <pre className="result">{helloResult}</pre>}
      </section>
    </div>
  )
}

export default App
