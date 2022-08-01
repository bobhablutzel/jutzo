<template>
  <image-map :image-map='this.imageMap()' src="/images/EA Framework.svg" class="imageMap" on="Business & Strategy"/>
  <div class="text-container">
    <p>The Business and Strategy discipline focuses on making sure that the company's strategic intent
    is will understood by the technologists, and that there is clarity on the priorities of the
    business so that investment decisions are appropriately aligned.</p>
    <p>A big part of the Business and Strategy discipline is developing (or customizing) a Business
    Capability model. Business capability models describe <i>what</i> the business does, without regard
    to <i>how</i> the business does it. They also provide the ability to categorize capabilities based
    on criteria such as criticality, centralization, and others.</p>
    <p>This discipline also gives an opportunity to discuss other business concerns, such as the
    organizational model. Ideally, this is tightly aligned to the capabilities (and, later, the
    platforms) but often diverges for a number of reasons. This divergence can be managed, but only
    if it is surfaced and understood.</p>
    <p>The Business and Strategy discipline is unique among the disciplines in that the Enterprise
      Architecture team is <router-link :to="{name: 'ARSCIV_X_Overview'}">responsible</router-link>
      for creating and maintaining the artifacts, but not for
    the thought leadership that underlies those artifacts. Put another way, while the Enterprise
    Architecture team often creates the Business Capability model, they are not in the position to
    define what Business Capabilities the company should provide.</p>
  </div>
</template>

<script>

import ImageMap from "@/components/ImageMap"
import frameworkImageMap from "@/views/framework/FrameworkImageMap"

export default {
  name: "DisciplineBusinessStrategy",
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