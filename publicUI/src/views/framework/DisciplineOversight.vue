<template>
  <image-map :image-map='this.imageMap()' src="/images/EA Framework.svg" class="imageMap" on="Oversight"/>
  <div class="text-container">
    <p>Oversight...</p>
  </div>
</template>

<script>

import ImageMap from "@/components/ImageMap"
import frameworkImageMap from "@/views/framework/FrameworkImageMap"

export default {
  name: "FiveDisciplines",
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