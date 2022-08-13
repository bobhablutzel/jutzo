<!---

Vue3 component to implement an image map. The image map projects
rectangular, clickable areas on top of an image.

The component has the following properties:

  src:            The source of the image
  alt:            The alt-text for the image
  zIndex:         The root z-index of the image (default: 10).
  on:             An optional string denoting the name of the currently active hot spot
  imageMap:       An object, consisting of
  |  height:      The height of the image map
  |  width:       The width of the image map
  |  hotspots:    An array of map records, where each map record contains
  |  |  coords:   An object containing the coordinates
  |  |  |  x:     The x location of the start of the hot spot (left)
  |  |  |  y:     The y location of the start of the hot spot (top)
  |  |  |  w:     The width of the hot spot
  |  |  |  h:     The height of the hot spot
  |  |  name:     An (optional) name
  |  |  to:       A vue-router link

--->

<template>
  <div ref="container" :style="cssVars" class="container" @mouseenter="resizeMap">
      <slot>
        <img :src="this.src" :alt="this.alt" class="image" ref="image"/>
      </slot>
    <canvas ref="canvas" class="imageMapCanvas"
            @mousemove="trackMouse($event)" @mouseleave="clearCurrentMap()"
            @click="navigate()" @resize="resizeMap($event)"/>
  </div>
</template>

<script>

export default {
  name: "ImageMap",
  props: {
    imageMap: Object,
    src: String,
    on: String,
    alt: { type: String, default: '' },
    zIndex: { type: Number, default: 10 }
  },

  data() {
    return {
      currentMap: null,         // The currently active map (might be null)
      scaleX: 0,                // The relative scale between the current image size
      scaleY: 0                 //    and the size that was used to create the image map
    }
  },



  mounted() {

    // Detect when the window size changes, so that we can react to the image size changing
    window.addEventListener('resize', this.resizeMap )
    this.draw()
  },

  unmounted() {
    window.removeEventListener('resize', this.resizeMap )
  },

  methods: {

    // Called when the image size changes; resizes the canvas and
    // recalculates the scale of the canvas to image map
    resizeMap() {
      const [ width, height ] = [ this.$refs.container.offsetWidth, this.$refs.container.offsetHeight ]
      if (this.$refs.canvas.width !== width || this.$refs.canvas.height !== height) {

        this.$refs.canvas.width = width;
        this.$refs.canvas.height = height;

        this.scaleX = width / this.imageMap.width
        this.scaleY = height / this.imageMap.height

        // If there is an active hot spot ($props.on), then
        // highlight that in a slightly darker shade
        if (this.on && this.on !== '') {

          // Get the fill color for an "on" field
          let fillColor = getComputedStyle(this.$refs.canvas).getPropertyValue('--image-map-on')
          if (!fillColor) fillColor = 'lightgreen'

          // Draw the on map
          const ctx = this.$refs.canvas.getContext("2d");
          this.imageMap.hotspots.forEach((hotspot) => {
            if (this.on === hotspot.name) {
              const coords = this.scaleCoords(hotspot.coords)
              ctx.save()
              ctx.fillStyle = fillColor
              ctx.globalAlpha = 0.2
              ctx.fillRect(coords.x, coords.y, coords.w, coords.h)
              ctx.restore()
            }
          })
        }
      }
    },

    // Called when the image map is clicked. This will navigate to
    // the right place if clicked in a hotspot
    navigate() {
      if (this.currentMap && this.currentMap.name !== this.on) {
        this.$router.push( this.currentMap.to )
      }
    },

    // Scales a set of coordinates to the current size of the image
    scaleCoords( coords ) {
      return {
        x: coords.x * this.scaleX,
        y: coords.y * this.scaleY,
        w: coords.w * this.scaleX,
        h: coords.h * this.scaleY
      }
    },

    // Tests the mouse location to see if it's in the rect
    pointInRect(mouseLoc, coords) {
      const localCoords = this.scaleCoords(coords)
      return (mouseLoc.x >= localCoords.x && mouseLoc.x <= (localCoords.x + localCoords.w) &&
          mouseLoc.y >= localCoords.y && mouseLoc.y <= (localCoords.y + localCoords.h))
    },

    // Looks for a map that corresponds to the mouse location. Sets
    // the current map if found
    attemptToFindCurrentMap(mouseLoc) {
      this.$props.imageMap.hotspots.forEach( (map) => {
        const coords = map.coords;
        if (this.pointInRect(mouseLoc, coords)) {
          if (this.currentMap !== map) {
            this.currentMap = map
            this.draw()
          }
        }
      } )
    },

    // Clears the current map (and redraws)
    clearCurrentMap() {
      if (this.currentMap) {
        const coords = this.scaleCoords(this.currentMap.coords)
        const ctx = this.$refs.canvas.getContext("2d");
        ctx.clearRect(coords.x, coords.y, coords.w, coords.h)

        if (this.currentMap.name === this.on) {
          // Get the fill color for an "on" field
          let fillColor = getComputedStyle(this.$refs.canvas).getPropertyValue('--image-map-on' )
          if (!fillColor) fillColor = 'lightgreen'

          ctx.save()
          ctx.fillStyle = fillColor
          ctx.globalAlpha = 0.2;
          ctx.fillRect(coords.x, coords.y, coords.w, coords.h)
          ctx.restore()
        }

      }
      this.currentMap = null
      this.draw()

    },


    // Draw the active map (if any) on top of the image
    draw() {
      this.resizeMap()
      const ctx = this.$refs.canvas.getContext("2d");

      // If we have a current map, draw the highlight
      if (this.currentMap && this.currentMap.name !== this.on) {

        // Get the fill color for a highlighted field
        let fillColor = getComputedStyle(this.$refs.canvas).getPropertyValue('--image-map-hover' )
        if (!fillColor) fillColor = 'lightgray'

        const coords = this.scaleCoords(this.currentMap.coords)
        ctx.save()
        ctx.fillStyle = fillColor
        ctx.globalAlpha = 0.2;
        ctx.fillRect(coords.x, coords.y, coords.w, coords.h)
        ctx.restore()
      }

    },

    trackMouse(evt) {
      this.resizeMap()

      // Get the mouse location relative to the canvas
      const rect = this.$refs.canvas.getBoundingClientRect();
      const mouseLoc = {
        x: (evt.clientX - rect.left),
        y: (evt.clientY - rect.top),
      }

      // If we have a current map, we check to see if we're still in the map
      if (this.currentMap) {

        const coords = this.currentMap.coords
        if (!this.pointInRect(mouseLoc, coords)) {

          // No longer in the current map, clear it and
          // see if we can find a new current map
          this.clearCurrentMap()
          this.attemptToFindCurrentMap(mouseLoc)
        }
      } else {

        // We didn't have a map. See if we do now.
        this.attemptToFindCurrentMap(mouseLoc)
      }

    },
  },

  computed: {
    cssVars() {
      return {
        '--base-z-index': this.zIndex,
      }
    },

    contents() {
      const slot = this.$slots.default ? this.$slots.default()[0].el : this.$refs.image
      console.log( slot )
      return slot
    }
  }

}

</script>


<style scoped>

/*noinspection CssUnresolvedCustomProperty*/
.container {
  display: grid;
  grid-auto-columns: auto;
  grid-auto-rows: auto;
  z-index: var(--base-z-index);
}


/*noinspection CssUnresolvedCustomProperty*/
.image {
  z-index: var(--base-z-index) + 1;
}

/*noinspection CssUnresolvedCustomProperty*/
.imageMapCanvas {
  position: absolute;
  z-index: var(--base-z-index) + 2;
}

</style>