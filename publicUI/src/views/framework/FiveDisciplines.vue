<template>
  <image-map :image-map='this.imageMap()' src="/images/EA Framework.svg" style="height: 300px; width: 800px;"/>
  <div class="text-container">
    <p>Architecture is a complex topic. Since architecture is
      <router-link :to="{ name: 'DefiningArchitecture'}">fundamentally about the non-functional aspects of a system,</router-link>
      any system of more than trivial complexity will start to have an architecture defined. The richness is valuable,
      because we do have to think about each of the architectures, but at the same time it creates confusion about
      which architectures are similar to each other, how they interact, and which facets of architecture can be
      reused across systems.
    </p>
    <p>
      The Scalable Enterprise Architecture Framework addresses this by grouping architecture into five major
      disciplines. These disciplines align different architectures by similarity of intent. This means that
      architectures that are more engineering focused are grouped together, and separately from those that are more
      product or business oriented. Additionally, the disciplines include not just the architectures but the
      disciplines involved in managing architecture: creating and communicating the expectations, and working to
      ensure that that guidance is followed.
    </p>
    <p>
      Arranging architecture in this way makes it simpler to understand how to manage architecture, because
      each of the disciplines needs different skills, stakeholders, and approaches.
    </p>
    <p>
      The Five Disciplines are:</p>
      <ul>
        <li><router-link :to="{name: 'DisciplineBusinessStrategy'}">Business and Strategy</router-link>: Capturing business intent</li>
        <li><router-link :to="{name: 'DisciplinePortfolioManagement'}">Portfolio Management</router-link>: Managing and mapping the portfolio of systems </li>
        <li><router-link :to="{name: 'DisciplineImplementation'}">Implementations</router-link>: Building sustainable systems</li>
        <li><router-link :to="{name: 'DisciplineGuidance'}">Guidance</router-link>: Communicating architectural constraints</li>
        <li><router-link :to="{name: 'DisciplineOversight'}">Oversight</router-link>: Evaluating compliance</li>
      </ul>
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

</style>