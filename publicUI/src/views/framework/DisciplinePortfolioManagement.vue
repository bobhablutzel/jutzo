<template>
  <image-map :image-map='this.imageMap()' src="/images/EA Framework.svg" class="imageMap" on="Portfolio Management"/>
  <div class="text-container">
    <p>Portfolio Management...</p>
  </div>
</template>

<script>

import ImageMap from "@/components/ImageMap"
import frameworkImageMap from "@/views/framework/FrameworkImageMap"

export default {
  name: "DisciplinePortfolioManagement",
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