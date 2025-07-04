import React, { useState } from 'react';
import { View, Image, StyleSheet, Dimensions, ActivityIndicator } from 'react-native';
import * as SvgComponent from 'react-native-svg';

type Point = {x: number, y: number};

type Props = {
  imageUri: string;
  boundingPoly?: Point[];
};

const SmartBackgroundRemoval: React.FC<Props> = ({ imageUri, boundingPoly }) => {
  const [isLoading, setIsLoading] = useState(true);
  
  // Si no hay datos de segmentación, usar el componente simple
  if (!boundingPoly || boundingPoly.length < 3) {
    return (
      <View style={styles.container}>
        <View style={styles.whiteBackground} />
        {isLoading && (
          <View style={styles.loaderContainer}>
            <ActivityIndicator size="small" color="#ee1e1e" />
          </View>
        )}
        <Image 
          source={{ uri: imageUri }} 
          style={styles.image} 
          resizeMode="contain"
          onLoadStart={() => setIsLoading(true)}
          onLoad={() => setIsLoading(false)}
          onLoadEnd={() => setIsLoading(false)}
        />
      </View>
    );
  }

  // Obtener dimensiones para el SVG
  const width = Dimensions.get('window').width - 32;
  const height = width;

  // Crear el path para el SVG
  const svgPath = boundingPoly
    .map((point, index) => {
      const x = point.x * width;
      const y = point.y * height;
      return `${index === 0 ? 'M' : 'L'}${x},${y}`;
    })
    .join(' ') + ' Z';

  return (
    <View style={styles.container}>
      <View style={styles.whiteBackground} />
      
      {isLoading && (
        <View style={styles.loaderContainer}>
          <ActivityIndicator size="small" color="#ee1e1e" />
        </View>
      )}
      
      {/* Imagen original */}
      <Image 
        source={{ uri: imageUri }} 
        style={styles.image}
        onLoadStart={() => setIsLoading(true)}
        onLoad={() => setIsLoading(false)}
        onLoadEnd={() => setIsLoading(false)}
      />
      
      {/* SVG con clip path */}
      <View style={StyleSheet.absoluteFill}>
        <SvgComponent.Svg height={height} width={width}>
          <SvgComponent.Defs>
            <SvgComponent.ClipPath id="clip">
              <SvgComponent.Path d={svgPath} />
            </SvgComponent.ClipPath>
          </SvgComponent.Defs>
          
          <SvgComponent.Rect
            x="0"
            y="0"
            width={width}
            height={height}
            fill="white"
            clipPath="url(#clip)"
          />
        </SvgComponent.Svg>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    width: '100%',
    aspectRatio: 1,
    borderRadius: 8,
    overflow: 'hidden',
    backgroundColor: 'white',
  },
  whiteBackground: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'white',
  },
  image: {
    width: '100%',
    height: '100%',
  },
  loaderContainer: {
    ...StyleSheet.absoluteFillObject,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(255,255,255,0.7)',
  }
});

export default SmartBackgroundRemoval;