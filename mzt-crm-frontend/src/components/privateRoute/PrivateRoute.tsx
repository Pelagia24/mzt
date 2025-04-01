import {Navigate, Outlet} from 'react-router-dom';
import { useRefreshTokenQuery } from '../../api/authApi.ts';

const PrivateRoute = () => {
    const { isLoading, isError } = useRefreshTokenQuery();

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (isError) {
        return <Navigate to='/login' />;
    }
    return <Outlet />;
};

export default PrivateRoute;