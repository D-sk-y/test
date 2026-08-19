import React, { useState } from 'react'
import '../css/All.css'

// 游戏结束条件
const WIN_COUNT = 5;

const DIRECTIONS = [
    [0, 1],   // 水平 →（同一行，列 +1）
    [1, 0],   // 垂直 ↓（行 +1，同一列）
    [1, 1],   // 主对角线 ↘（行 +1，列 +1）
    [1, -1],  // 副对角线 ↙（行 +1，列 -1）
];

// 棋子标记：0 = 空，1 = 黑，2 = 白
const EMPTY = 0;
const BLACK = 1;
const WHITE = 2;

function Home() {
    // history：每一步的棋面快照（二维数组的数组）
    const [history, setHistory] = useState<number[][][]>([
        Array.from({ length: 9 }, () => Array(9).fill(EMPTY)),
    ]);
    // 步数
    const [step, setStep] = useState<number>(0);
    // 当前棋面 = history[step]
    const board = history[step];

    // 先手
    const [player, setPlayer] = useState<number>(BLACK);
    // 赢家
    const [winner, setWinner] = useState<number | null>(null);
    // 每一步的描述
    const moves = history.map((_, move) => (
        <li key={move}>
            <button onClick={() => jumpTo(move)}>
                {move === 0 ? 'Go to game start' : 'Go to move #' + move}</button>
        </li>
    ));

    // 处理点击事件
    function handleClick(row: number, col: number) {
        // 该位置已有棋子，或已经分出胜负，忽略这次点击
        if (board[row][col] !== EMPTY || winner !== null) {
            return;
        }
        // copy 
        const newBoard = board.map((r) => [...r]);
        newBoard[row][col] = player;
        //  history
        const nextHistory = [...history.slice(0, step + 1), newBoard];
        setHistory(nextHistory);
        setStep(nextHistory.length - 1);

        console.log(row, col);
        console.log(history);

        // 判断这步是否获胜
        if (checkWin(newBoard, row, col)) {
            setWinner(player);
        } else {
            // 没赢就换对手下
            setPlayer(player === BLACK ? WHITE : BLACK);
        }
    }

    function reset() {
        setHistory([Array.from({ length: 9 }, () => Array(9).fill(EMPTY))]);
        setStep(0);
        setPlayer(BLACK);
        setWinner(null);
    }

    function jumpTo(move: number) {
        setStep(move);
    }

    return (
        <div className="HomeBackground" >
            <Board board={board} handleClick={handleClick} reset={reset} moves={moves} />
        </div>
    )
}

// 棋盘上的一个格子
function Square({ value, onSquareClick }) {
    return (
        <button className="square" onClick={onSquareClick}>
            {value === BLACK ? '黑' : value === WHITE ? '白' : ''}
        </button>
    )
}

// 棋盘
function Board({ board, handleClick, reset, moves }) {
    return (
        <div className="board">
            {Array.from({ length: 9 }, (_, row) => (
                <div className="row" key={row}>
                    {Array.from({ length: 9 }, (_, col) => (
                        <Square key={row * 9 + col} value={board[row][col]} onSquareClick={() => handleClick(row, col)} />
                    ))}
                </div>
            ))}
            <div style={{ textAlign: 'center', marginTop: 12 }}>
                <button onClick={reset} style={{ padding: '6px 24px', fontSize: 16, cursor: 'pointer' }}>
                    重新开始
                </button>
            </div>
            <div>
                <ol>{moves}</ol>
            </div>
        </div>
    )
}

// 从 (row, col) 出发，沿 (dr, dc) 方向数有几个连续的同色棋子
function countInDirection(board: number[][], row: number, col: number, dr: number, dc: number, player: number): number {
    let count = 0;
    let r = row + dr;
    let c = col + dc;

    // 一直在棋盘内，且还是同一颜色的棋子，就继续数
    while (r >= 0 && r < 9 && c >= 0 && c < 9 && board[r][c] === player) {
        count++;
        r += dr;
        c += dc;
    }
    return count;
}

// 每次落子后调用：返回 true 表示赢了
function checkWin(board: number[][], row: number, col: number): boolean {
    // 刚落的棋子是谁的
    const player = board[row][col];
    // 空位不可能赢
    if (player === 0)
        return false;

    for (const [dr, dc] of DIRECTIONS) {
        // 朝正方向数 + 朝反方向数，加上自己 = 该方向上的连续棋子数
        const count =
            1 +
            countInDirection(board, row, col, dr, dc, player) +
            countInDirection(board, row, col, -dr, -dc, player);

        // 任意方向连成 5 颗，赢了
        if (count >= 5)
            return true;
    }
    return false;
}

export default Home