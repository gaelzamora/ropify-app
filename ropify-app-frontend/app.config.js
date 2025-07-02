export default {
  expo: {
    name: "ropify-app-frontend",
    slug: "ropify-app-frontend",
    scheme: "ropify",
    version: "1.0.0",
    orientation: "portrait",
    icon: "./assets/images/icon.png",
    userInterfaceStyle: "automatic",
    android: {
      package: "com.gaelzamora.ropifyappfrontend",
      adaptiveIcon: {
        foregroundImage: "./assets/images/adaptive-icon.png",
        backgroundColor: "#ffffff"
      }
    },
    ios: {
      supportsTablet: true
    },
    web: {
      bundler: "metro",
      favicon: "./assets/images/favicon.png"
    },
    extra: {
      BASE_URL: "http://192.168.1.78:8080",
      CLIENT_IOS_ID: "581346763419-c7j70hh1q75djb3n6jglipgor9nk57f4.apps.googleusercontent.com",
      CLIENT_WEB_ID: "581346763419-0b4visk3isdvfsrvtqal75qpsdmug0ie.apps.googleusercontent.com",
      CLIENT_ANDROID_ID: "581346763419-14krm7ggp05u07p7v364jcdqofj0vo11.apps.googleusercontent.com"
    }
  }
};