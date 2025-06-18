import React, { useState, useEffect } from 'react';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

function App() {
    const [userId, setUserId] = useState('user1');
    const [amount, setAmount] = useState(100);
    const [description, setDescription] = useState('Test order');
    const [balance, setBalance] = useState(0);
    const [depositAmount, setDepositAmount] = useState(50);
    const [orders, setOrders] = useState([]);

    const fetchData = async (url, method = 'GET', body = null) => {
        try {
            const options = {
                method,
                headers: { 'Content-Type': 'application/json' },
            };
            if (body) options.body = JSON.stringify(body);

            const response = await fetch(url, options);
            if (!response.ok) {
                const error = await response.text();
                throw new Error(error);
            }
            return await response.json();
        } catch (error) {
            toast.error(error.message);
            throw error;
        }
    };

    const createAccount = async () => {
        await fetchData('http://localhost:8081/accounts', 'POST', { user_id: userId });
        toast.success('Account created successfully!');
    };

    const deposit = async () => {
        await fetchData(`http://localhost:8081/accounts/${userId}/deposit`, 'POST', { amount: parseFloat(depositAmount) });
        const { balance } = await fetchData(`http://localhost:8081/accounts/${userId}/balance`);
        setBalance(balance);
        toast.success('Deposit successful!');
    };

    const createOrder = async () => {
        const order = await fetchData('http://localhost:8080/orders', 'POST', {
            user_id: userId,
            amount: parseFloat(amount),
            description,
        });
        setOrders([...orders, order]);
        toast.success(`Order ${order.id} created!`);
    };

    const loadOrders = async () => {
        const orders = await fetchData('http://localhost:8080/orders');
        setOrders(orders);
    };

    const loadBalance = async () => {
        try {
            const { balance } = await fetchData(`http://localhost:8081/accounts/${userId}/balance`);
            setBalance(balance);
        } catch (error) {
            setBalance(0);
        }
    };

    useEffect(() => {
        loadBalance();
        loadOrders();
    }, []);

    return (
        <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
            <h1>Online Shop</h1>

            <div style={{ marginBottom: '20px', padding: '15px', border: '1px solid #ddd', borderRadius: '5px' }}>
                <h2>Account Management</h2>
                <div style={{ marginBottom: '10px' }}>
                    <label>User ID: </label>
                    <input
                        value={userId}
                        onChange={(e) => setUserId(e.target.value)}
                        style={{ marginRight: '10px' }}
                    />
                    <button onClick={createAccount}>Create Account</button>
                </div>
                <div style={{ marginBottom: '10px' }}>
                    <label>Deposit Amount: </label>
                    <input
                        type="number"
                        value={depositAmount}
                        onChange={(e) => setDepositAmount(e.target.value)}
                        style={{ marginRight: '10px', width: '100px' }}
                    />
                    <button onClick={deposit} style={{ marginRight: '10px' }}>Deposit</button>
                    <span>Balance: {balance.toFixed(2)}</span>
                </div>
            </div>

            <div style={{ marginBottom: '20px', padding: '15px', border: '1px solid #ddd', borderRadius: '5px' }}>
                <h2>Order Management</h2>
                <div style={{ marginBottom: '10px' }}>
                    <label>Amount: </label>
                    <input
                        type="number"
                        value={amount}
                        onChange={(e) => setAmount(e.target.value)}
                        style={{ marginRight: '10px', width: '100px' }}
                    />
                </div>
                <div style={{ marginBottom: '10px' }}>
                    <label>Description: </label>
                    <input
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        style={{ marginRight: '10px', width: '300px' }}
                    />
                </div>
                <button onClick={createOrder} style={{ marginRight: '10px' }}>Create Order</button>
                <button onClick={loadOrders}>Refresh Orders</button>
            </div>

            <div style={{ padding: '15px', border: '1px solid #ddd', borderRadius: '5px' }}>
                <h2>Order List</h2>
                {orders.length === 0 ? (
                    <p>No orders yet</p>
                ) : (
                    <ul>
                        {orders.map((order) => (
                            <li key={order.id} style={{ marginBottom: '10px' }}>
                                <strong>ID:</strong> {order.id} <br />
                                <strong>Amount:</strong> {order.amount} <br />
                                <strong>Status:</strong> {order.status} <br />
                                <strong>Description:</strong> {order.description}
                            </li>
                        ))}
                    </ul>
                )}
            </div>

            <ToastContainer position="bottom-right" autoClose={5000} />
        </div>
    );
}

export default App;