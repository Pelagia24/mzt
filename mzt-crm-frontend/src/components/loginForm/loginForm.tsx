import {useEffect, useState} from 'react';
import {Button} from 'primereact/button';
import {InputText} from 'primereact/inputtext';
import {classNames} from 'primereact/utils';
import authFormStyles from './loginForm.module.scss';
import {NavLink, useNavigate} from "react-router-dom";
import {useLoginUserMutation} from "../../api/authApi.ts";


export const LoginForm = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [submitted, setSubmitted] = useState(false);
    const [login, {isSuccess}] = useLoginUserMutation();
    const navigate = useNavigate();

    useEffect(() => {
        if (isSuccess) {
            navigate('/');
        }
    }, [isSuccess]);

    const isEmailValid = (email: string) => /^\S+@\S+\.\S+$/.test(email);
    const isPasswordValid = (password: string) => password.length > 8;

    return (
        <div className={authFormStyles.container}>
            <div className="surface-card p-4 shadow-2 border-round w-30rem">
                <div className="text-center mb-5">
                    <img src="/logo.png" alt="Логотип MZT" height={100} className="mb-3"/>
                    <div className="text-900 text-3xl font-medium mb-3">Добро пожаловать!</div>
                    <span className="text-600 font-medium line-height-3">Нет аккаунта?</span>
                    <NavLink to={'/signup'}>&nbsp;Зарегистрируйтесь</NavLink>
                </div>

                <div className="field">
                    <label className="block text-900 font-medium mb-2">Email</label>
                    <InputText
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className={classNames({'p-invalid': submitted && (!email || !isEmailValid(email))}, 'w-full mb-1')}
                    />
                    {submitted && !email && <small className="p-error">Введите email</small>}
                    {submitted && email && !isEmailValid(email) &&
                        <small className="p-error">Неверный email</small>}
                </div>

                <div className="field">
                    <label htmlFor="password" className="block text-900 font-medium mb-2">Пароль</label>
                    <InputText
                        id="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        type="password"
                        className={classNames({'p-invalid': submitted && (!password || !isPasswordValid(password))}, 'w-full mb-1')}
                    />
                    {submitted && !password && <small className="p-error">Введите пароль</small>}
                    {submitted && password && !isPasswordValid(password) &&
                        <small className="p-error">Мин. 8 символов</small>}
                </div>

                <Button
                    label="Войти"
                    icon="pi pi-user"
                    className="w-full mt-1"
                    onClick={() => {
                        setSubmitted(true);
                        if (!email || !password || !isEmailValid(email) || !isPasswordValid(password)) {
                            return;
                        }
                        login({
                            email,
                            password
                        })
                    }}
                />
            </div>
        </div>
    );
};

export default LoginForm;
