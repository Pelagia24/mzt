export interface RegisterUserRequest  {
    email: string;
    password: string;
    age: number;
    birthdate: string;
    city: string;
    employment: string;
    is_business_owner: string;
    month_income: number;
    name: string;
    phone_number: string;
    position_at_work: string;
    telegram?: string
}

export interface LoginUserRequest {
    email: string;
    password: string;
}