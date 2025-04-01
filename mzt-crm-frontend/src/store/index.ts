import { configureStore } from '@reduxjs/toolkit';
import authApi from '../api/authApi.ts';
import authSlice from "./slices/authSlice.ts";

export const store = configureStore({
    reducer: {
        [authApi.reducerPath]: authApi.reducer,
        authSlice
    },
    middleware: getDefaultMiddleware => getDefaultMiddleware().concat(authApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export default store;
