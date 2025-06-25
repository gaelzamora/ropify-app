import React, { useState, useEffect } from "react";
import { View, Text, TextInput, TouchableOpacity, StyleSheet, KeyboardAvoidingView, Platform, ActivityIndicator } from "react-native";
import { FontAwesome } from '@expo/vector-icons';
import { useRouter, useLocalSearchParams } from "expo-router";
import { useAuth } from "@/context/AuthContext";
import * as WebBrowser from 'expo-web-browser';
import * as Linking from 'expo-linking';
import AsyncStorage from '@react-native-async-storage/async-storage';
import * as SecureStore from 'expo-secure-store';
import * as Google from 'expo-auth-session/providers/google'
import { Image } from "react-native";

// Registrar para recibir el resultado de la autenticación
WebBrowser.maybeCompleteAuthSession();

if (typeof window !== "undefined") {
  window.onload = () => {
    if (window.location.hash.includes("access_token")) {
      console.log("YA ENTRE (onload)");
      const params = new URLSearchParams(window.location.hash.substring(1));
      const accessToken = params.get("access_token");
      if (accessToken) {
        // Puedes usar un evento personalizado o guardar el token en localStorage
        window.localStorage.setItem("google_access_token", accessToken);
        // Limpia el hash
        window.location.hash = "";
        // Opcional: recarga la app o navega a la ruta deseada
        window.location.reload();
      }
    }
  };
}

const Login: React.FC = () => {
    const { authenticate, setIsLoggedIn, setUser, isLoadingAuth } = useAuth();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isGoogleLoading, setIsGoogleLoading] = useState(false);
    const [manualToken, setManualToken] = useState("")
    const router = useRouter();

    const [request, response, promptAsync] = Google.useAuthRequest({
        iosClientId: '581346763419-1v79br463i6ufa9mks5jd1kasj74nkr3.apps.googleusercontent.com',
        webClientId: '581346763419-0b4visk3isdvfsrvtqal75qpsdmug0ie.apps.googleusercontent.com',
        scopes: ['profile', 'email'],
        redirectUri: 'https://auth.expo.io/@gaelzamora/ropify-app-frontend',
        usePKCE: false
    })

    console.log("Aqui ando")

    // Manejar deep links entrantes (cuando vuelve de la autenticación)
    useEffect(() => {
      if (typeof window !== "undefined") {
        const accessToken = window.localStorage.getItem("google_access_token");
        if (accessToken) {
            console.log("Access token capturado desde localStorage:", accessToken);
            handlerGoogleAuthentication(accessToken);
            window.localStorage.removeItem("google_access_token");
        } else {
            console.log("AQUI NO HAY NADA ")
        }
    }  
    }, []);
    
    async function handlerGoogleAuthentication(accessToken: string | undefined) {
        console.log("HELOOOO")

        if (!accessToken) {
            console.log("No entre")
            return;
        } 
        
        try {
            console.log("Entre aqui")

            setIsGoogleLoading(true);
            console.log("Access token recibido:", accessToken.substring(0, 10) + "...");
            
            // URL de tu backend (usa tu IP real)
            const baseUrl = 'http://192.168.1.79:8080';
            const response = await fetch(`${baseUrl}/api/oauth/google/token`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ access_token: accessToken })
            });
            
            const data = await response.json();
            console.log("Respuesta del backend:", data);
            
            // Guardar token y datos de usuario
            if (data.status === 'success' && data.data?.token) {
                await SecureStore.setItemAsync('token', data.data.token); 
                await AsyncStorage.setItem('user', JSON.stringify(data.data.user));
                
                setUser(data.data.user);
                setIsLoggedIn(true);
                router.replace('/(authed)/(tabs)/settings');
            } else {
                console.error("Error en respuesta:", data);
            }
        } catch (error) {
            console.error('Error en la autenticación con Google:', error);
        } finally {
            setIsGoogleLoading(false);
        }
    }

    const handlerGoogleLogin = () => {
        promptAsync()
    }
    
    // Autenticación normal
    async function onAuthenticate() {
        await authenticate("login", email, password);
    }
    

    return (
        <KeyboardAvoidingView
            style={styles.container}
            behavior={Platform.OS === "ios" ? "padding" : undefined}
        >
            <TouchableOpacity style={styles.backButton} onPress={() => router.push("/")}>
                <FontAwesome name="arrow-left" size={25} color="#211" />
            </TouchableOpacity>
            <View style={{ alignItems: "center" }}>
                <Image
                    source={require("@/assets/images/logo/letras.png")}
                    style={{ width: "50%", height: 50 }} 
                />
            </View>
            <Text style={styles.title}>Sign in</Text>

            <View style={styles.form}>
                {/* Campos de login existentes */}
                <TextInput
                    style={styles.input}
                    placeholder="Email"
                    keyboardType="email-address"
                    autoCapitalize="none"
                    autoCorrect={false}
                    placeholderTextColor={"#888"}
                    cursorColor={"#888"}
                    value={email}
                    onChangeText={setEmail}
                />
                <TextInput
                    style={styles.input}
                    placeholder="Password"
                    secureTextEntry
                    autoCapitalize="none"
                    autoCorrect={false}
                    placeholderTextColor={"#888"}
                    cursorColor={"#888"}
                    value={password}
                    onChangeText={setPassword}
                />
                <TouchableOpacity 
                    style={styles.button}
                    onPress={onAuthenticate}
                    disabled={isLoadingAuth}
                >
                    {isLoadingAuth ? (
                        <ActivityIndicator color="#fff" />
                    ) : (
                        <Text style={styles.buttonText}>Login</Text>
                    )}
                </TouchableOpacity>
                <TouchableOpacity 
                    onPress={() => router.push('/register')}
                    style={styles.loginLink}
                >
                    <Text style={styles.loginText}>Already not have an account? <Text style={styles.loginTextBold}>Sign Up</Text></Text>
                </TouchableOpacity>
            </View>
            
            <Text style={{ fontSize: 12, color: "#888", textAlign: "center", marginTop: 20 }}>Or sign in with</Text>
            
            <View style={styles.socialContainer}>
                <TouchableOpacity 
                    style={styles.socialButton}
                    onPress={handlerGoogleLogin}
                    disabled={isGoogleLoading}
                >
                    {isGoogleLoading ? (
                        <ActivityIndicator size="small" color="#DB4437" />
                    ) : (
                        <FontAwesome name="google" size={28} color="#DB4437" />
                    )}
                </TouchableOpacity>
                <TouchableOpacity style={styles.socialButton}>
                    <FontAwesome name="facebook" size={28} color="#1877F3" />
                </TouchableOpacity>
                <TouchableOpacity style={styles.socialButton}>
                    <FontAwesome name="twitter" size={28} color="#1DA1F2" />
                </TouchableOpacity>
            </View>
            
            <Text style={{marginTop: 20, color: "#888", fontSize: 12}}>
                ¿Problemas con el login de Google? Pega aquí el access_token:
            </Text>
            <TextInput
                style={styles.input}
                placeholder="Pega el access_token aquí"
                value={manualToken}
                onChangeText={setManualToken}
            />
            <TouchableOpacity
                style={styles.button}
                onPress={() => console.log("HEY APRETANDO")}
            >
                <Text style={styles.buttonText}>Continuar con token</Text>
            </TouchableOpacity>

        </KeyboardAvoidingView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: "#fff",
        justifyContent: "center",
        alignItems: "center",
        padding: 16,
        position: "relative"
    },
    backButton: {
        position: "absolute",
        top: 60,
        left: 26,
        zIndex: 10,
        backgroundColor: "transparent",
        padding: 4,
    },
    form: {
        width: "100%",
        maxWidth: 340,
        borderRadius: 8,
        padding: 24,
        shadowColor: "#000",
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.08,
        shadowRadius: 8,
        elevation: 2,
    },
    title: {
        fontSize: 28,
        fontWeight: "600",
        marginBottom: 24,
        color: "#222",
        marginTop: 15
    },
    input: {
        height: 48,
        borderColor: "#ddd",
        borderWidth: 1,
        borderRadius: 6,
        paddingHorizontal: 12,
        marginBottom: 16,
        backgroundColor: "#fff",
        fontSize: 16,
    },
    button: {
        backgroundColor: "#222",
        paddingVertical: 14,
        borderRadius: 6,
        alignItems: "center",
        marginTop: 8,
    },
    buttonText: {
        color: "#fff",
        fontWeight: "600",
        fontSize: 16,
    },
    socialContainer: {
        flexDirection: "row",
        justifyContent: "center",
        marginTop: 24,
        gap: 24,
    },
    socialButton: {
        padding: 8,
    },
    loginLink: {
        marginTop: 20,
        alignItems: "center",
    },
    loginText: {
        color: "#555",
        fontSize: 14,
    },
    loginTextBold: {
        fontWeight: "600",
        color: "#222",
    },
});

export default Login;