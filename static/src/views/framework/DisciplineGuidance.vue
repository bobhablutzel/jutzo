<template>
  <image-map :image-map='this.imageMap()' src="/images/EA Framework.svg" class="imageMap" on="Guidance"/>
  <div class="text-container">
    <p>Guidance...</p>
  </div>
</template>

<script>

import ImageMap from "@/components/ImageMap"
import frameworkImageMap from "@/views/framework/FrameworkImageMap"

export default {
  name: "DisciplineGuidance",
  components: { ImageMap },
  methods: {
    imageMap() { return frameworkImageMap }
  }
}


</script>

<style scoped>


.imageMap {
  height: 300px;
  width: 800px;
}
</style>