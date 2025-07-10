import React, { useState, useEffect, useCallback } from 'react';
import { View, Image, StyleSheet, Dimensions } from 'react-native';
import * as SvgComponent from 'react-native-svg';
import Shimmer from './Shimmer';

type Point = {x: number, y: number};

type Props = {
  imageUri: string;
  boundingPoly?: Point[];
  style?: any;
};

const SmartBackgroundRemoval: React.FC<Props> = ({ imageUri, boundingPoly, style }) => {
  const [isLoading, setIsLoading] = useState(true);
  
  useEffect(() => {
    let timeoutId: NodeJS.Timeout;
    
    if (imageUri) {
      timeoutId = setTimeout(() => {
        setIsLoading(false);
      }, 10000);
    }
    
    return () => {
      if (timeoutId) clearTimeout(timeoutId);
    };
  }, [imageUri]);
  
  const handleLoadStart = useCallback(() => {
    setIsLoading(true);
  }, []);
  
  const handleLoadSuccess = useCallback(() => {
    setIsLoading(false);
  }, []);
  
  const handleLoadEnd = useCallback(() => {
    setTimeout(() => {
      setIsLoading(false);
    }, 100);
  }, []);
  
  if (!boundingPoly || boundingPoly.length < 3) {
    return (
      <View style={[styles.container, style]}>
        <View style={styles.whiteBackground} />
        
        {/* Mostrar el shimmer solo si isLoading es true */}
        {isLoading && (
          <View style={styles.loaderContainer}>
            <Shimmer />
          </View>
        )}
        
        <Image 
          source={{ uri: imageUri }} 
          style={[
            styles.image, 
            isLoading ? { opacity: 0.4 } : { opacity: 1 }
          ]} 
          resizeMode="contain"
          onLoadStart={handleLoadStart}
          onLoad={handleLoadSuccess}
          onLoadEnd={handleLoadEnd}
          // Añadir esta prop para priorizar la carga
          progressiveRenderingEnabled={true}
        />
      </View>
    );
  }

  // El resto del código para el caso con boundingPoly
  const width = Dimensions.get('window').width - 32;
  const height = width;
  
  const svgPath = boundingPoly
    .map((point, index) => {
      const x = point.x * width;
      const y = point.y * height;
      return `${index === 0 ? 'M' : 'L'}${x},${y}`;
    })
    .join(' ') + ' Z';

  return (
    <View style={[styles.container, style]}>
      <View style={styles.whiteBackground} />
      
      {isLoading && (
        <View style={styles.loaderContainer}>
          <Shimmer />
        </View>
      )}
      
      <Image 
        source={{ uri: imageUri }} 
        style={[styles.image, isLoading ? { opacity: 0.4 } : { opacity: 1 }]}
        onLoadStart={handleLoadStart}
        onLoad={handleLoadSuccess}
        onLoadEnd={handleLoadEnd}
        progressiveRenderingEnabled={true}
      />
      
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
    backgroundColor: 'transparent',
    zIndex: 5,
  }
});

export default SmartBackgroundRemoval;