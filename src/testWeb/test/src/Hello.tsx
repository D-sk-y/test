import React from 'react'
import './clo.css'
import { useState } from 'react'

interface HelloProps {
  name: string
}

const Hello: React.FC<HelloProps> = ({ name }) => {
  const [count, setCount] = useState(0)
  return (
    <div>
      <h1>Hello, {name}!</h1>
      <button className="btn" onClick={() => {
        setCount(count + 1)
        console.log(count)
      }}>Click me</button>
    </div>
  )
}

export default Hello