// components/SmartBackgroundRemoval.tsx
import React from 'react';
import { View, Image, StyleSheet, Dimensions } from 'react-native';
import Svg, { Path, Defs, ClipPath, Rect } from 'react-native-svg';

type Point = {x: number, y: number};

type Props = {
  imageUri: string;
  boundingPoly?: Point[];
};

const SmartBackgroundRemoval: React.FC<Props> = ({ imageUri, boundingPoly }) => {
  // Si no hay datos de segmentación, usar el componente simple
  if (!boundingPoly || boundingPoly.length < 3) {
    return (
      <View style={styles.container}>
        <View style={styles.whiteBackground} />
        <Image source={{ uri: imageUri }} style={styles.image} resizeMode="contain" />
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
      
      {/* Imagen original */}
      <Image source={{ uri: imageUri }} style={styles.image} />
      
      {/* SVG con clip path */}
      <View style={StyleSheet.absoluteFill}>
        <Svg height={height} width={width}>
          <Defs>
            <ClipPath id="clip">
              <Path d={svgPath} />
            </ClipPath>
          </Defs>
          
          <Rect
            x="0"
            y="0"
            width={width}
            height={height}
            fill="white"
            clipPath="url(#clip)"
          />
        </Svg>
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
});

export default SmartBackgroundRemoval;