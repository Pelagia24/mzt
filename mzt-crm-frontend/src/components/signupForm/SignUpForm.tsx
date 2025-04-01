import {useEffect, useState} from 'react';
import {InputText} from 'primereact/inputtext';
import {Calendar} from 'primereact/calendar';
import {Dropdown} from 'primereact/dropdown';
import {Button} from 'primereact/button';
import {classNames} from 'primereact/utils';
import authFormStyles from './signUpForm.module.scss';
import {NavLink, useNavigate} from "react-router-dom";
import {Nullable} from "primereact/ts-helpers";
import {useRegisterUserMutation} from "../../api/authApi.ts";

export const SignUpForm = () => {
    const [registerUser, {isSuccess}] = useRegisterUserMutation();
    const navigate = useNavigate();
    const [step, setStep] = useState(1);
    const [formData, setFormData] = useState({
        fullName: '',
        birthDate: null,
        email: '',
        password: '',
        phone: '',
        city: '',
        age: '',
        profession: '',
        isOwner: '',
        position: '',
        income: '',
        telegram: '',
    });
    const [submitted, setSubmitted] = useState(false);

    useEffect(() => {
        if (isSuccess) {
            navigate('/');
        }
    }, [isSuccess]);

    const isEmailValid = (email: string) => /^\S+@\S+\.\S+$/.test(email);
    const isPasswordValid = (password: string) => password.length > 8;
    const isPhoneValid = (phone: string) => /^\+?\d{10,15}$/.test(phone);

    const handleChange = (field: string, value: string | Nullable<Date>) => {
        setFormData({...formData, [field]: value});
    };

    const validateStep1 = () =>
        formData.fullName &&
        formData.birthDate &&
        formData.email &&
        isEmailValid(formData.email) &&
        formData.password &&
        isPasswordValid(formData.password);

    const validateStep2 = () =>
        formData.phone &&
        isPhoneValid(formData.phone) &&
        formData.city &&
        formData.age &&
        formData.profession;

    const validateStep3 = () =>
        formData.isOwner &&
        formData.position &&
        formData.income;

    const onSubmit = () => {
        setSubmitted(true);
        if (validateStep3()) {
            registerUser({
                email: formData.email,
                password: formData.password,
                age: parseInt(formData.age),
                birthdate: (formData.birthDate as unknown as Date).toISOString(),
                city: formData.city,
                employment: formData.profession,
                is_business_owner: formData.isOwner,
                month_income: parseInt(formData.income),
                name: formData.fullName,
                position_at_work: formData.position,
                phone_number: formData.phone,
                telegram: formData.telegram
            })
        }
    };

    return (
        <div className={authFormStyles.container}>
            <div className="surface-card p-4 shadow-2 border-round w-30rem">
                <div className="text-center mb-5">
                    <img src="/logo.png" alt="Логотип MZT" height={100} className="mb-3"/>
                    <div className="text-900 text-3xl font-medium mb-3">Давайте знакомиться!</div>
                    <span className="text-600 font-medium line-height-3">Есть аккаунт?</span>
                    <NavLink to={'/login'}>&nbsp;Войдите</NavLink>
                </div>
                {step === 1 && (
                    <>
                        <div className="field">
                            <label className="block text-900 font-medium mb-2">ФИО</label>
                            <InputText
                                value={formData.fullName}
                                onChange={(e) => handleChange('fullName', e.target.value)}
                                className={classNames({'p-invalid': submitted && !formData.fullName}, 'w-full mb-1')}
                            />
                            {submitted && !formData.fullName && <small className="p-error">Введите ФИО</small>}
                        </div>

                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Дата рождения</label>
                            <Calendar
                                value={formData.birthDate}
                                onChange={(e) => handleChange('birthDate', e.value)}
                                showIcon
                                dateFormat="dd.mm.yy"
                                className={classNames({'p-invalid': submitted && !formData.birthDate}, 'w-full mb-1')}
                            />
                            {submitted && !formData.birthDate &&
                                <small className="p-error">Выберите дату рождения</small>}
                        </div>

                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Email</label>
                            <InputText
                                value={formData.email}
                                onChange={(e) => handleChange('email', e.target.value)}
                                className={classNames({'p-invalid': submitted && (!formData.email || !isEmailValid(formData.email))}, 'w-full mb-1')}
                            />
                            {submitted && !formData.email && <small className="p-error">Введите email</small>}
                            {submitted && formData.email && !isEmailValid(formData.email) &&
                                <small className="p-error">Неверный email</small>}
                        </div>

                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Пароль</label>
                            <InputText
                                type="password"
                                value={formData.password}
                                onChange={(e) => handleChange('password', e.target.value)}
                                className={classNames({'p-invalid': submitted && (!formData.password || !isPasswordValid(formData.password))}, 'w-full mb-1')}
                            />
                            {submitted && !formData.password && <small className="p-error">Введите пароль</small>}
                            {submitted && formData.password && !isPasswordValid(formData.password) &&
                                <small className="p-error">Мин. 8 символов</small>}
                        </div>

                        <Button label="Далее" onClick={() => {
                            setSubmitted(true);
                            if (validateStep1()) {
                                setStep(2);
                                setSubmitted(false);
                            }
                        }}/>
                    </>
                )}

                {step === 2 && (
                    <>
                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Номер телефона</label>
                            <InputText
                                value={formData.phone}
                                onChange={(e) => handleChange('phone', e.target.value)}
                                className={classNames({'p-invalid': submitted && (!formData.phone || !isPhoneValid(formData.phone))}, 'w-full mb-1')}
                            />
                            {submitted && !formData.phone && <small className="p-error">Введите номер телефона</small>}
                            {submitted && formData.phone && !isPhoneValid(formData.phone) &&
                                <small className="p-error">Неверный формат</small>}
                        </div>

                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Город</label>
                            <InputText
                                value={formData.city}
                                onChange={(e) => handleChange('city', e.target.value)}
                                className={classNames({'p-invalid': submitted && !formData.city}, 'w-full mb-1')}
                            />
                            {submitted && !formData.city && <small className="p-error">Введите город проживания</small>}
                        </div>

                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Возраст</label>
                            <InputText
                                value={formData.age}
                                onChange={(e) => handleChange('age', e.target.value)}
                                className={classNames({'p-invalid': submitted && !formData.age}, 'w-full mb-1')}
                            />
                            {submitted && !formData.age && <small className="p-error">Введите свой возраст</small>}
                        </div>

                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Сфера деятельности</label>
                            <InputText
                                value={formData.profession}
                                onChange={(e) => handleChange('profession', e.target.value)}
                                className={classNames({'p-invalid': submitted && !formData.profession}, 'w-full mb-1')}
                            />
                            {submitted && !formData.profession &&
                                <small className="p-error">Введите вашу сферу деятельности</small>}
                        </div>

                        <Button label="Назад" className="p-button-secondary mr-2" onClick={() => setStep(1)}/>
                        <Button label="Далее" onClick={() => {
                            setSubmitted(true);
                            if (validateStep2()) {
                                setStep(3);
                                setSubmitted(false);
                            }
                        }}/>
                    </>
                )}

                {step === 3 && (
                    <>
                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Вы собственник бизнеса?</label>
                            <Dropdown
                                value={formData.isOwner}
                                options={[
                                    {label: 'Да', value: 'yes'},
                                    {label: 'Нет', value: 'no'},
                                    {label: 'Другое', value: 'other'}
                                ]}
                                onChange={(e) => handleChange('isOwner', e.value)}
                                placeholder="Выберите"
                                className={classNames({'p-invalid': submitted && !formData.isOwner}, 'w-full mb-1')}
                            />
                            {submitted && !formData.isOwner && <small className="p-error">Выберите вариант</small>}
                        </div>

                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Должность</label>
                            <InputText
                                value={formData.position}
                                onChange={(e) => handleChange('position', e.target.value)}
                                className={classNames({'p-invalid': submitted && !formData.position}, 'w-full mb-1')}
                            />
                            {submitted && !formData.position &&
                                <small className="p-error">Введите свою должность</small>}
                        </div>

                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Ежемесячный личный доход</label>
                            <InputText
                                value={formData.income}
                                onChange={(e) => handleChange('income', e.target.value)}
                                className={classNames({'p-invalid': submitted && !formData.income}, 'w-full mb-1')}
                            />
                            {submitted && !formData.income && <small className="p-error">Введите Ваш доход</small>}
                        </div>

                        <div className="field">
                            <label className="block text-900 font-medium mb-2">Ник в Телеграме (необязательно)</label>
                            <InputText
                                value={formData.telegram}
                                onChange={(e) => handleChange('telegram', e.target.value)}
                                className={'w-full mb-1'}
                            />
                        </div>

                        <Button label="Назад" className="p-button-secondary mr-2" onClick={() => setStep(2)}/>
                        <Button label="Зарегистрироваться" onClick={onSubmit}/>
                    </>
                )}
            </div>
        </div>
    );
};
