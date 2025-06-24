import React, { useState, useEffect } from "react";
import { View, Text, TextInput, TouchableOpacity, StyleSheet, KeyboardAvoidingView, Platform, ActivityIndicator } from "react-native";
import { FontAwesome } from '@expo/vector-icons';
import { useRouter, useLocalSearchParams } from "expo-router";
import { useAuth } from "@/context/AuthContext";
import * as WebBrowser from 'expo-web-browser';
import * as Linking from 'expo-linking';
import AsyncStorage from '@react-native-async-storage/async-storage';
import * as SecureStore from 'expo-secure-store';

// Registrar para recibir el resultado de la autenticación
WebBrowser.maybeCompleteAuthSession();

const Login: React.FC = () => {
    const { authenticate, setIsLoggedIn, setUser, isLoadingAuth } = useAuth();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isGoogleLoading, setIsGoogleLoading] = useState(false);
    const router = useRouter();
    const params = useLocalSearchParams();
    
    // Manejar deep links entrantes (cuando vuelve de la autenticación)
    useEffect(() => {
        if (params && params.data) {
            handleDeepLink(params.data as string);
        }
        
        // Suscribirse a deep links futuros
        const subscription = Linking.addEventListener('url', ({url}) => {
            const { queryParams } = Linking.parse(url);
            if (queryParams && queryParams.data) {
                handleDeepLink(queryParams.data as string);
            }
        });
        
        return () => {
            subscription.remove();
        };
    }, [params]);
    
    // Procesar datos de autenticación recibidos por deep link
    const handleDeepLink = async (encodedData: string) => {
        try {
            const decoded = decodeURIComponent(encodedData);
            const jsonData = Buffer.from(decoded, 'base64').toString();
            const authData = JSON.parse(jsonData);
            
            if (authData.token) {
                await SecureStore.setItemAsync('userToken', authData.token);
                
                // Obtener detalles del usuario usando el token
                const response = await fetch('http://192.168.1.79:8080/api/auth/me', {
                    headers: { 
                        'Authorization': `Bearer ${authData.token}`
                    }
                });
                const userData = await response.json();
                
                setUser(userData.data);
                setIsLoggedIn(true);
                await AsyncStorage.setItem('user', JSON.stringify(userData.data));
                
                router.replace('/(authed)/(tabs)/settings');
            }
        } catch (error) {
            console.error('Error procesando datos de autenticación:', error);
        }
    };
    
    // Autenticación normal
    async function onAuthenticate() {
        await authenticate("login", email, password);
    }
    
    // Iniciar autenticación de Google
    async function handleGoogleLogin() {
        try {
            setIsGoogleLoading(true);
            
            // URL de redirección a nuestra app
            const redirectUri = Linking.createURL('oauth-callback');
            
            // Abrir navegador con la URL de login de Google
            const baseUrl = Platform.OS === 'android' 
                ? 'http://192.168.1.XXX:8080' 
                : 'http://localhost:8080';
                
            const result = await WebBrowser.openAuthSessionAsync(
                `${baseUrl}/api/oauth/google/login?redirect_uri=${encodeURIComponent(redirectUri)}`,
                redirectUri
            );
            
            if (result.type === 'success') {
                // La autenticación exitosa será manejada por el efecto que escucha los deep links
                console.log("Login exitoso, esperando datos...");
            }
        } catch (error) {
            console.error("Error de autenticación con Google:", error);
        } finally {
            setIsGoogleLoading(false);
        }
    }

    return (
        <KeyboardAvoidingView
            style={styles.container}
            behavior={Platform.OS === "ios" ? "padding" : undefined}
        >
            {/* Contenido existente */}
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
                    onPress={handleGoogleLogin}
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