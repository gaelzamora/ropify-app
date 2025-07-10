import React, { useRef, useEffect } from 'react';
import { View, StyleSheet, Animated } from 'react-native';
import { LinearGradient } from 'expo-linear-gradient';

type ShimmerProps = {
  width?: number | string;
  height?: number | string;
  style?: any;
}

const Shimmer: React.FC<ShimmerProps> = ({ width = '100%', height = '100%', style }) => {
  const animatedValue = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    Animated.loop(
      Animated.timing(animatedValue, {
        toValue: 1,
        duration: 3000, // Aumentar la duración para que sea más lenta (de 1500ms a 3000ms)
        useNativeDriver: true,
      })
    ).start();
  }, []);

  // Cambiamos a translateY para moverlo verticalmente
  const translateY = animatedValue.interpolate({
    inputRange: [0, 1],
    outputRange: [-350, 350], // Mantener los mismos valores pero ahora para movimiento vertical
  });

  return (
    <View style={[styles.container, { width, height }, style]}>
      <Animated.View
        style={[StyleSheet.absoluteFill, { transform: [{ translateY }] }]} // Cambiar translateX por translateY
      >
        <LinearGradient
          style={StyleSheet.absoluteFill}
          // Cambiar dirección del gradiente de horizontal a vertical
          start={{ x: 0, y: 0 }} 
          end={{ x: 0, y: 1 }}   // Cambiado de {x: 1, y: 0} a {x: 0, y: 1}
          colors={['#f0f0f0', '#ffffff', '#f0f0f0']}
          locations={[0.3, 0.5, 0.7]}
        />
      </Animated.View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: '#f0f0f0',
    overflow: 'hidden',
    borderRadius: 8,
  },
});

export default Shimmer; 