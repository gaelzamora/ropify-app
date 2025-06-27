import { userService } from "@/services/user";
import AsyncStorage from "@react-native-async-storage/async-storage";
import React, { useEffect, useState } from "react";
import { User } from "@/types/user";
import { router } from "expo-router";
import { oauthService } from "@/services/oauth";

interface AuthContextProps {
    isLoggedIn: boolean;
    isLoadingAuth: boolean;
    authenticate: (authMode: "login" | "register", email: string, password: string, username?: string, firstName?: string, lastName?: string) => Promise<void>
    logout: VoidFunction;
    user: User | null;
    authenticateWithGoogle: (accessToken: string | undefined) => Promise<void>;
    setIsLoggedIn: (state: boolean) => void;
    setUser: (user: User) => void
}

const AuthContext = React.createContext({} as AuthContextProps);

export function useAuth() {
    return React.useContext(AuthContext)
}

export function AuthenticationProvider({ children }: React.PropsWithChildren) {
    const [isLoggedIn, setIsLoggedIn] = useState(false)
    const [isLoadingAuth, setIsLoadingAuth] = useState(false)
    const [user, setUser] = useState<User | null>(null)

    useEffect(() => {
        async function checkIfLoggedIn() {
            const token = await AsyncStorage.getItem("token")
            const user = await AsyncStorage.getItem("user")

            if (token && user) {
                setIsLoggedIn(true)
                setUser(JSON.parse(user))
                router.replace("/(authed)/(tabs)/settings" as any)
            } else {
                setIsLoggedIn(false)
            }
        }

        checkIfLoggedIn()
    }, [])

    async function authenticate(
        authMode: "login" | "register",  
        email: string, 
        password: string, 
        username?: string, 
        firstName?: string, 
        lastName?: string
    ) {
        try {
            setIsLoadingAuth(true)

            if (authMode === "register") {
                // Validar campos obligatorios
                if (!username || !username.trim()) {
                    setIsLoadingAuth(false)
                    return
                }
                if (!firstName || !firstName.trim() || !lastName || !lastName.trim()) {
                    setIsLoadingAuth(false)
                    return
                }
            }

            let response;
            if (authMode === "login") {
                response = await userService[authMode](email, password);
            } else if (authMode === "register" && firstName && lastName && username) {
                response = await userService[authMode](firstName, lastName, username, email, password);
            }

            if (response) {
                    setIsLoggedIn(true)
                    await AsyncStorage.setItem("token", response.data.token)
                    await AsyncStorage.setItem("user", JSON.stringify(response.data.user))
                    setUser(response.data.user)
                    router.replace("/(authed)/(tabs)/(closet)" as any)
                
            }
        } catch (error) {
            console.log(error)
            setIsLoggedIn(false)
        } finally {
            setIsLoadingAuth(false)
        }
    }

    async function authenticateWithGoogle(accessToken: string | undefined) {
        try {
            setIsLoadingAuth(true);
            if (accessToken) {
                const response = await oauthService.getTokenFromGoogle(accessToken);

                if (response) {
                    setIsLoggedIn(true);
                    await AsyncStorage.setItem("token", response.data.token);
                    await AsyncStorage.setItem("user", JSON.stringify(response.data.user));
                    setUser(response.data.user);
                    router.replace("/(authed)/(tabs)/(closet)" as any);
                    console.log("New user: ", user)
                } else {
                    console.log("No hay response")
                }
                
            } else {
                console.log("No hay access token que mandar.")
            }
        } catch (error) {
            console.error("Google authentication failed:", error);
            // Mostrar un toast o alerta al usuario
        } finally {
            setIsLoadingAuth(false);
        }
    }

    async function logout() {
        setIsLoggedIn(false)
        await AsyncStorage.removeItem("token")
        await AsyncStorage.removeItem("user")
    }

    return (
        <AuthContext.Provider
            value={ {
                authenticate,
                logout,
                isLoggedIn,
                isLoadingAuth,
                user,
                authenticateWithGoogle,
                setIsLoggedIn,
                setUser
            } }>
                { children }
        </AuthContext.Provider>
    )
}