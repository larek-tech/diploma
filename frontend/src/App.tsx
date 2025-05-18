import {useEffect} from 'react';
import {Route, Routes} from 'react-router-dom';
import {AuthProvider} from './auth/AuthProvider';
import {RequireAuth} from './auth/RequireAuth';
import {RequireUnauth} from './auth/RequireUnauth';
import {Dashboard} from './components/Dashboard';
import HomeBanner from './components/HomeBanner';
import {Toaster} from './components/ui/toaster';
import Chat from './pages/Chat';
import ChatFromDomain from './pages/ChatFromDomain';
import ChatFromHistory from './pages/ChatFromHistory';
import CreateDomain from './pages/CreateDomain';
import {Login} from './pages/Login';
import {Pages} from './router/constants';

function App() {
    useEffect(() => {
        // eslint-disable-next-line
        // @ts-ignore
        if (window.Telegram && window.Telegram.WebApp) {
            // eslint-disable-next-line
            // @ts-ignore
            const telegramWebApp = window.Telegram.WebApp;

            telegramWebApp.expand();
        } else {
            console.error('Telegram WebApp is not available');
        }
    }, []);

    return (
        <>
            <Toaster />

            <AuthProvider>
                <Routes>
                    <Route
                        path={`/${Pages.Login}`}
                        element={
                            <RequireUnauth>
                                <Login />
                            </RequireUnauth>
                        }
                    />
                    <Route
                        path={`/${Pages.Chat}/:sessionId`}
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <ChatFromHistory />
                                </Dashboard>
                            </RequireAuth>
                        }
                    />
                    <Route
                        path={`/${Pages.Domain}/:domainId`}
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <ChatFromDomain />
                                </Dashboard>
                            </RequireAuth>
                        }
                    />
                    <Route
                        path='*'
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <HomeBanner />
                                </Dashboard>
                            </RequireAuth>
                        }
                    />
                    <Route
                        path={`/${Pages.CreateDomain}`}
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <CreateDomain />
                                </Dashboard>
                            </RequireAuth>
                        }
                    />
                </Routes>
            </AuthProvider>
        </>
    );
}

export default App;
