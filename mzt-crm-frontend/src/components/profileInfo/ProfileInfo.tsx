import {Button} from "primereact/button";
import {classNames} from "primereact/utils";
import styles from "./ProfileInfo.module.scss";
import {useGetProfileQuery} from "../../api/authApi.ts";

const ProfileInfo = () => {
    const {data, isLoading, isError} = useGetProfileQuery();


    if (isLoading) {
        return <div className={classNames("text-center surface-card border-round m-4 p-4 shadow-2", styles.wrapper)}>
            Загрузка...
        </div>
    }

    if (isError) {
        return <div className={classNames("text-center surface-card border-round m-4 p-4 shadow-2", styles.wrapper)}>
            Произошла ошибка, попробуйте позже
        </div>
    }

    return data ? (<div className={classNames("surface-card border-round m-4 p-4 shadow-2", styles.wrapper)}>
        <div className="font-medium text-3xl text-900 mb-3">Моя информация</div>
        <ul className="list-none p-0 m-0">
            <li className="flex align-items-center py-3 px-2 border-top-1  flex-wrap">
                <div className="text-500 w-6 md:w-2 font-bold">ФИО</div>
                <div className="text-900 w-full md:w-8 md:flex-order-0 flex-order-1">{data!.user.name}</div>
                <div className="w-6 md:w-2 flex justify-content-end">
                    <Button label="Изменить" icon="pi pi-pencil" className="p-button-text"/>
                </div>
            </li>
            <li className="flex align-items-center py-3 px-2 border-top-1  flex-wrap">
                <div className="text-500 w-6 md:w-2 font-bold">Дата рождения</div>
                <div className="text-900 w-full md:w-8 md:flex-order-0 flex-order-1">
                    {new Date(data!.user.birthdate).toLocaleDateString()}
                </div>
                <div className="w-6 md:w-2 flex justify-content-end">
                    <Button label="Изменить" icon="pi pi-pencil" className="p-button-text"/>
                </div>
            </li>
            <li className="flex align-items-center py-3 px-2 border-top-1  flex-wrap">
                <div className="text-500 w-6 md:w-2 font-bold">Номер телефона</div>
                <div className="text-900 w-full md:w-8 md:flex-order-0 flex-order-1">{data!.user.phone_number}</div>
                <div className="w-6 md:w-2 flex justify-content-end">
                    <Button label="Изменить" icon="pi pi-pencil" className="p-button-text"/>
                </div>
            </li>
            <li className="flex align-items-center py-3 px-2 border-top-1  flex-wrap">
                <div className="text-500 w-6 md:w-2 font-bold">Город</div>
                <div className="text-900 w-full md:w-8 md:flex-order-0 flex-order-1">{data!.user.city}</div>
                <div className="w-6 md:w-2 flex justify-content-end">
                    <Button label="Изменить" icon="pi pi-pencil" className="p-button-text"/>
                </div>
            </li>
            <li className="flex align-items-center py-3 px-2 border-top-1  flex-wrap">
                <div className="text-500 w-6 md:w-2 font-bold">Сфера деятельности</div>
                <div className="text-900 w-full md:w-8 md:flex-order-0 flex-order-1 line-height-3">
                    {data!.user.employment}
                </div>
                <div className="w-6 md:w-2 flex justify-content-end">
                    <Button label="Изменить" icon="pi pi-pencil" className="p-button-text"/>
                </div>
            </li>
            <li className="flex align-items-center py-3 px-2 border-top-1  flex-wrap">
                <div className="text-500 w-6 md:w-2 font-bold">Должность</div>
                <div className="text-900 w-full md:w-8 md:flex-order-0 flex-order-1 line-height-3">
                    {data!.user.position_at_work}
                </div>
                <div className="w-6 md:w-2 flex justify-content-end">
                    <Button label="Изменить" icon="pi pi-pencil" className="p-button-text"/>
                </div>
            </li>
            <li className="flex align-items-center py-3 px-2 border-top-1  flex-wrap">
                <div className="text-500 w-6 md:w-2 font-bold">Ежемесячный доход</div>
                <div className="text-900 w-full md:w-8 md:flex-order-0 flex-order-1 line-height-3">
                    {`${data!.user.month_income}₽`}
                </div>
                <div className="w-6 md:w-2 flex justify-content-end">
                    <Button label="Изменить" icon="pi pi-pencil" className="p-button-text"/>
                </div>
            </li>
            <li className="flex align-items-center py-3 px-2 border-top-1 border-bottom-1  flex-wrap">
                <div className="text-500 w-6 md:w-2 font-bold">Telegram</div>
                <div className="text-900 w-full md:w-8 md:flex-order-0 flex-order-1 line-height-3">
                    {data!.user.telegram || 'Не указан'}
                </div>
                <div className="w-6 md:w-2 flex justify-content-end">
                    <Button label="Изменить" icon="pi pi-pencil" className="p-button-text"/>
                </div>
            </li>
        </ul>
    </div>) : null;
}

export default ProfileInfo;