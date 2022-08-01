<template>
  <panel header="Scalable Enterprise Architecture Framework">
    <p>There are two major components of the Scalable Enterprise Architecture Framework:</p>
    <div class="panel-container">
      <card>
        <template #title>
          <i class="pi pi-pf pi-sitemap" style="padding-top: 5px; padding-left: 5px;"></i>
          Discipline Map
        </template>
        <template #content>
          <image-map :image-map='this.imageMap()' src="/images/EA Framework.svg"/>
          <span class="caption">The Scalable Enterprise Architecture Discipline Map</span>
          <ul>
            <li>Defines <router-link :to="{ name: 'FrameworkFiveDisciplines'}">the five major architecture disciplines</router-link> and their function;</li>
            <li>Creates a natural staffing model for architecture offices;</li>
            <li>Scales to any size organization</li>
          </ul>
        </template>
      </card>
      <divider layout="vertical"/>
      <card>
        <template #title>
          <i class="pi pi-pf pi-sync" style="padding-top: 5px; padding-left: 5px;"></i>
          Operational Procedures
        </template>
        <template #content>
          <ul>
            <li>Prescribes how architects at all levels of the organization interact</li>
            <li>Eliminates bottlenecks and creates efficiencies</li>
            <li>Ensures proper reviews and documentation occur naturally</li>
          </ul>
        </template>
      </card>
    </div>
    <p>
      In order to understand the framework,
    </p>
  </panel>

</template>


<script>

import ImageMap from "@/components/ImageMap"
import frameworkImageMap from "@/views/framework/FrameworkImageMap"


export default {
  name: "FrameworkOverview",
  components: { ImageMap },
  methods: {
    imageMap() { return frameworkImageMap }
  }
}
</script>

<style scoped>

</style>