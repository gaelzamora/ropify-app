import React, { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, StyleSheet, KeyboardAvoidingView, Platform, Image, ScrollView, Alert, ActivityIndicator } from "react-native";
import { FontAwesome } from '@expo/vector-icons';
import { useRouter } from "expo-router";
import { useAuth } from "@/context/AuthContext";

const Register: React.FC = () => {
    const router = useRouter();

    const { authenticate, isLoadingAuth } = useAuth()
    const [firstName, setFirstName] = useState("")
    const [lastName, setLastName] = useState("")
    const [username, setUsername] = useState("")
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [repassword, setRePassword] = useState("")

    async function onAuthenticate() {
        if (password === repassword) await authenticate("register", email, password, username, firstName, lastName)
        else Alert.alert("The passwod must be same password")
    }

    return (
        <KeyboardAvoidingView
            style={styles.container}
            behavior={Platform.OS === "ios" ? "padding" : undefined}
        >
            <TouchableOpacity style={styles.backButton} onPress={() => router.push("/")}>
                <FontAwesome name="arrow-left" size={25} color="#211" />
            </TouchableOpacity>

            <ScrollView 
                contentContainerStyle={styles.scrollContainer} 
                showsVerticalScrollIndicator={false}
            >
                <View style={styles.form}>
                    <View style={{ alignItems: "center" }}>
                        <Image
                            source={require("@/assets/images/logo/letras.png")}
                            style={{ width: "50%", height: 50 }}
                        />
                    </View>
                    <Text style={styles.title}>Create Account</Text>
                    
                    <View style={styles.nameContainer}>
                        <TextInput
                            style={[styles.input, styles.halfInput]}
                            placeholder="First Name"
                            autoCapitalize="words"
                            autoCorrect={false}
                            placeholderTextColor={"#888"}
                            cursorColor={"#888"}
                            value={firstName}
                            onChangeText={setFirstName}
                        />
                        <TextInput
                            style={[styles.input, styles.halfInput]}
                            placeholder="Last Name"
                            autoCapitalize="words"
                            autoCorrect={false}
                            placeholderTextColor={"#888"}
                            cursorColor={"#888"}
                            value={lastName}
                            onChangeText={setLastName}
                        />
                    </View>

                    <View style={styles.usernameContainer}>
                        <Text style={styles.atSymbol}>@</Text>
                        <TextInput
                            style={styles.usernameInput}
                            placeholder="Username"
                            autoCapitalize="none"
                            autoCorrect={false}
                            placeholderTextColor={"#888"}
                            cursorColor={"#888"}
                            value={username}
                            onChangeText={setUsername}
                        />
                    </View>
                    
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
                    <TextInput
                        style={styles.input}
                        placeholder="Confirm Password"
                        secureTextEntry
                        autoCapitalize="none"
                        autoCorrect={false}
                        placeholderTextColor={"#888"}
                        cursorColor={"#888"}
                        value={repassword}
                        onChangeText={setRePassword}
                    />
                    
                    <TouchableOpacity 
                        style={styles.button}
                        onPress={onAuthenticate}
                        disabled={isLoadingAuth}
                    >
                        {isLoadingAuth ? (
                            <ActivityIndicator color="#fff" />
                        ) : (
                            <Text style={styles.buttonText}>Create Account</Text>
                        )}
                    </TouchableOpacity>
                    
                    <TouchableOpacity 
                        onPress={() => router.push('/login')}
                        style={styles.loginLink}
                    >
                        <Text style={styles.loginText}>Already have an account? <Text style={styles.loginTextBold}>Sign In</Text></Text>
                    </TouchableOpacity>
                </View>
                
                <Text style={{ fontSize: 12, color: "#888", textAlign: "center", marginTop: 20 }}>Or sign up with</Text>
                <View style={styles.socialContainer}>
                    <TouchableOpacity style={styles.socialButton}>
                        <FontAwesome name="google" size={28} color="#DB4437" />
                    </TouchableOpacity>
                    <TouchableOpacity style={styles.socialButton}>
                        <FontAwesome name="facebook" size={28} color="#1877F3" />
                    </TouchableOpacity>
                    <TouchableOpacity style={styles.socialButton}>
                        <FontAwesome name="twitter" size={28} color="#1DA1F2" />
                    </TouchableOpacity>
                </View>
            </ScrollView>
        </KeyboardAvoidingView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: "#fff",
    },
    scrollContainer: {
        flexGrow: 1,
        alignItems: "center",
        justifyContent: "center",
        padding: 16,
        paddingTop: 80,
        paddingBottom: 40,
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
    nameContainer: {
        flexDirection: "row",
        justifyContent: "space-between",
    },
    halfInput: {
        width: "48%",
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
    socialContainer: {
        flexDirection: "row",
        justifyContent: "center",
        marginTop: 24,
        gap: 24,
    },
    socialButton: {
        padding: 8,
    },
    usernameContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        borderColor: "#ddd",
        borderWidth: 1,
        borderRadius: 6,
        height: 48,
        marginBottom: 16,
        paddingHorizontal: 12,
        backgroundColor: '#fff',
    },
    atSymbol: {
        fontSize: 16,
        color: "#888",
        marginRight: 8,
        fontWeight: '600',
    },
    usernameInput: {
        flex: 1,
        height: '100%',
        fontSize: 16,
        padding: 0,
    },
});

export default Register;