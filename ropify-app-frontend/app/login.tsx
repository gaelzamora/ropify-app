import React, { useState, useEffect } from "react";
import { View, Text, TextInput, TouchableOpacity, StyleSheet, KeyboardAvoidingView, Platform, ActivityIndicator, Image } from "react-native";
import { FontAwesome } from '@expo/vector-icons';
import { useRouter } from "expo-router";
import * as WebBrowser from "expo-web-browser";
import { useAuth } from "@/context/AuthContext";
import * as Google from 'expo-auth-session/providers/google'
import Constants from 'expo-constants'
import AsyncStorage from "@react-native-async-storage/async-storage";

WebBrowser.maybeCompleteAuthSession();


const Login: React.FC = () => {
    const { authenticate, isLoadingAuth, authenticateWithGoogle } = useAuth();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isGoogleLoading, setIsGoogleLoading] = useState(false);
    const router = useRouter();
    const [userInfo, setUserInfo] = useState<any>();

    const [request, response, promptAsync] = Google.useAuthRequest({
        iosClientId: Constants.expoConfig?.extra?.CLIENT_IOS_ID,
        webClientId: Constants.expoConfig?.extra?.CLIENT_WEB_ID,
        androidClientId: Constants.expoConfig?.extra?.CLIENT_ANDROID_ID,
    })

    useEffect(() => {
        console.log("Google Auth Response: ", response)
        handleEffect();
    }, [response]);

    async function handleEffect() {
        const user = await getLocalUser();
        console.log(user)
        if (!user) {
        if (response?.type === "success") {
            await authenticateWithGoogle(response?.authentication?.accessToken)
            //getUserInfo(response?.authentication?.accessToken ?? null);
        }
        } else {
            setUserInfo(user);
        }
    }
    const getLocalUser = async () => {
        const data = await AsyncStorage.getItem("@user");
        if (!data) return null;
        return JSON.parse(data);
    };

    const getUserInfo = async (token: string | null) => {
        if (!token) return;
        try {
            const response = await fetch(
                "https://www.googleapis.com/userinfo/v2/me",
                {
                headers: { Authorization: `Bearer ${token}` },
                }
            );

            const user = await response.json();
            setUserInfo(user);
        } catch (error) {
            console.log(error)
        }
    };

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
            {userInfo && <Text>{userInfo?.name}</Text>}
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
                    onPress={() => promptAsync()}
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