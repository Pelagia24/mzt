import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface AuthState {
    accessToken: string | null;
}

const initialState: AuthState = {
    accessToken: null,
}

interface CredentialsPayload {
    accessToken: string;
}

export const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setCredentials(state, action: PayloadAction<CredentialsPayload>) {
            state.accessToken = action.payload.accessToken;
        }
    }
});

export const {setCredentials} = authSlice.actions;
export default authSlice.reducer;