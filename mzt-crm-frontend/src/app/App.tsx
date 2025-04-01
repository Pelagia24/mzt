import {PrimeReactProvider} from 'primereact/api';
import {Route, BrowserRouter as Router, Routes} from 'react-router-dom';
import {LoginPage} from '../pages/loginPage/loginPage';
import {RegistrationPage} from '../pages/registrationPage/registrationPage';
import MainPage from "../pages/mainPage/MainPage.tsx";
import ProfilePage from "../pages/profilePage/ProfilePage.tsx";
import PrivateRoute from "../components/privateRoute/PrivateRoute.tsx";

function App() {
    return (
        <PrimeReactProvider>
            <Router>
                <Routes>
                    <Route path='/login' element={<LoginPage/>}/>
                    <Route path='/signup' element={<RegistrationPage/>}/>
                    <Route element={<PrivateRoute />}>
                        <Route path='/profile' element={<ProfilePage/>} />
                    </Route>
                    <Route path='/' element={<MainPage/>}/>
                </Routes>
            </Router>
        </PrimeReactProvider>
    );
}

export default App;
