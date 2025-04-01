import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {AuthResponse, UserInfoResponse} from "../types/authApi/responses.ts";
import {RegisterUserRequest, LoginUserRequest} from "../types/authApi/requests.ts";
import { setCredentials } from "../store/slices/authSlice.ts";
import {Dispatch} from "@reduxjs/toolkit";
import {RootState} from "../store";


const updateAccess = async (
    queryFulfilled: Promise<{ data: AuthResponse }>,
    dispatch: Dispatch,
) => {
    try {
        const {data} = await queryFulfilled;
        dispatch(setCredentials({accessToken: data.access_token}));
    } catch (error) {
        console.log(error);
    }
};

const authApi = createApi({
    reducerPath: 'api',
    tagTypes: ['Auth', 'UserInfo'],
    baseQuery: fetchBaseQuery({
        baseUrl: 'http://localhost:8080/api/v1',
        credentials: 'include',
        prepareHeaders: (headers, {getState, endpoint}) => {
            if (endpoint === 'refreshToken' || endpoint === 'getProfile') {
                const token = (getState() as RootState).authSlice.accessToken;
                if (token) {
                    headers.set('Authorization', `Bearer ${token}`);
                }
            }
            return headers;
        }
    }),
    endpoints: builder => ({
        registerUser: builder.mutation<AuthResponse, RegisterUserRequest>({
            query: body => ({
                url: '/auth/signup',
                method: 'POST',
                body,
            }),
            async onQueryStarted(_, {dispatch, queryFulfilled}) {
                await updateAccess(queryFulfilled, dispatch)
            },
            invalidatesTags: ['Auth'],
        }),
        loginUser: builder.mutation<AuthResponse, LoginUserRequest>({
            query: credentials => ({
                url: '/auth/signin',
                method: 'POST',
                body: credentials,
            }),
            async onQueryStarted(_, {dispatch, queryFulfilled}) {
                await updateAccess(queryFulfilled, dispatch)
            },
            invalidatesTags: ['Auth'],

        }),
        refreshToken: builder.query<AuthResponse, void>({
            query: () => ({
                url: '/auth/refresh',
                method: 'POST'
            }),
            providesTags: ['Auth'],
            async onQueryStarted(_, {dispatch, queryFulfilled}) {
                await updateAccess(queryFulfilled, dispatch)
            },
        }),
        getProfile: builder.query<UserInfoResponse, void>({
            query: () => '/users/me',
            providesTags: ['UserInfo'],
        })
    }),
});

export const {
    useRegisterUserMutation,
    useLoginUserMutation,
    useRefreshTokenQuery,
    useGetProfileQuery,
} = authApi;

export default authApi;
