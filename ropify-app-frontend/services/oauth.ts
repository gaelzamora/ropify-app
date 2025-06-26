import { Api } from "./api";

async function getTokenFromGoogle(access_token: string) {
  return Api.post("/oauth/google/token", { access_token })
}

const oauthService = {
  getTokenFromGoogle
}

export { oauthService }

// import * as WebBrowser from 'expo-web-browser';
// import * as AuthSession from 'expo-auth-session';
// import { Platform } from 'react-native';
// import { Api } from './api';
// import { AuthResponse } from '@/types/user';
// import Constants from 'expo-constants';

// // Registrar el esquema para el manejo de redirección
// WebBrowser.maybeCompleteAuthSession();

// // Configurar URLs de autenticación según plataforma
// const redirectUri = AuthSession.makeRedirectUri({
//   scheme: 'ropify',
//   path: 'oauth-callback'
// });

// // URL base de la API según plataforma
// const baseUrl = Platform.OS === 'android' 
//   ? Constants.expoConfig?.extra?.BASE_URL 
//   : Constants.expoConfig?.extra?.BASE_URL;

// const googleLogin = async (): Promise<AuthResponse> => {
//   try {
//     // 1. Iniciar el flujo de autenticación abriendo el navegador
//     const result = await WebBrowser.openAuthSessionAsync(
//       `${baseUrl}/api/oauth/google/login?redirect_uri=${encodeURIComponent(redirectUri)}`,
//       redirectUri
//     );

//     // 2. Procesar el resultado del navegador
//     if (result.type === 'success') {
//       // Extraer el código de autorización de la URL
//       const url = result.url;
//       const params = new URLSearchParams(url.split('?')[1]);
//       const code = params.get('code');
      
//       if (!code) {
//         throw new Error('No authorization code found');
//       }
      
//       // 3. Intercambiar el código por un token JWT usando nuestro backend
//       return await Api.post('/oauth/google/token', { code });
//     } else {
//       throw new Error('Authentication canceled or failed');
//     }
//   } catch (error) {
//     console.error('OAuth error:', error);
//     throw error;
//   }
// };

// export const oauthService = {
//   googleLogin
// };