<template>
  <image-map :image-map='this.imageMap()' src="/images/EA Framework.svg" style="height: 300px; width: 800px;"/>
  <div class="text-container">
    <p>Grouping the practice of architecture into
      <router-link :to="{name: 'FrameworkFiveDisciplines'}">five disciplines</router-link>
      does more than simply group like things together. It also provides a good model for creating an
      Enterprise Architecture office. Using this model, the Chief Architect would have key leaders in the
      organization for of the disciplines.
    </p>
    <p>
      Of course, that doesn't mean that you have to have five leads at the top level. Depending on the
      organization, some of the disciplines might be less active or performed by an external group.
      For example, many organizations have audit departments that are responsible for evaluating whether
      standard policies and procedures are being executed. In this world, the <i>Oversight</i> role
      might be more one of coordination with an audit team, and could be combined in an Operations
      role with the <i>Guidance</i> discipline. </p>
    <image-map alt="An example organization chart" :image-map='this.org1ImageMap()'>
      <svg class="orgChart">
        <path d="M250 1 h 220 v 60 h -220 z"/>
        <path d="M 360 60 v 10 M 90 70 H 610 v 20 M 90 70 v 20 M 263.33 70 v 20 M 436.66 70 v 20"/>
        <path d="M 20 90 h 140 v 40 h -140 z M 193.33 90 h 140 v 40 h -140 z"/>
        <path d="M 366.66 90 h 140 v 40 h -140 z M 540 90 h 140 v 40 h -140 z"/>
        <path d="M 610 130 v 10 M 570 140 H 650 v 15 M 570 140 v 15"/>
        <path d="M 538 155 h 65 v 30 h -65 z M 620 155 h 65 v 30 h -65 z"/>
        <text class="orgText level1" x="280" y="38" >Chief Architect</text>
        <text class="orgText level2" x="30" y="115">Business Architecture</text>
        <text class="orgText level2" x="200.33" y="115">Portfolio Management</text>
        <text class="orgText level2" x="374.66" y="115">Technical Architecture</text>
        <text class="orgText level2" x="576" y="115">Operations</text>
        <text class="orgText level3" x="546" y="174">Guidance</text>
        <text class="orgText level3" x="629" y="174">Oversight</text>
      </svg>
    </image-map>
    <p>
      Alternatively, you might have an organization that is largely focusing on the technology
      aspects of architecture and still maturing toward business alignment. That doesn't mean you
      should ignore the need to articulate business capabilities, but you might not put as much
      energy into that. In that case the <i>Business & Strategy</i> discipline might be rolled under
      and combined with the <i>Portfolio Management</i> discipline.</p>
    <image-map alt="Another example organization chart" :image-map='this.org2ImageMap()'>
      <svg class="orgChart">
        <path d="M250 1 h 220 v 60 h -220 z"/>
        <path d="M 360 60 v 10 M 142 70 H 610 v 20 m -158 -20 v 20 m -158 -20 v 20 m -151 -20 v 20"/>
        <path d="M 72 90 h 140 v 40 h -140 z M 228.33 90 h 140 v 40 h -140 z"/>
        <path d="M 384.66 90 h 140 v 40 h -140 z M 540 90 h 140 v 40 h -140 z"/>
        <path d="M 142 130 v 10 M 72.5 140 H 214.5 v 15 m -142 -15 v 15"/>
        <path d="M 10 155 h 125 v 30 h -125 z M 152 155 h 125 v 30 h -125 z"/>
        <text class="orgText level1"  x="280" y="38" >Chief Architect</text>
        <text class="orgText level2" x="85" y="115">Asset Management</text>
        <text class="orgText level2" x="236.33" y="115">Technical Architecture</text>
        <text class="orgText level2" x="395" y="115">Standards & Policies</text>
        <text class="orgText level2" x="555" y="115">Audit & Compliance</text>
        <text class="orgText level3" x="19" y="174">Business Capabilities</text>
        <text class="orgText level3" x="160" y="174">Portfolio Management</text>
      </svg>
    </image-map>

  </div>
</template>

<script>

import ImageMap from "@/components/ImageMap"
import frameworkImageMap from "@/views/framework/FrameworkImageMap"

export default {
  name: "OrganizationalAlignment",
  components: {ImageMap},
  methods: {
    imageMap() {
      return frameworkImageMap
    },
    org1ImageMap() {
      return {
        height: 200,
        width: 1000,
        hotspots: [
          {coords: {x: 250, y: 1, w: 220, h: 60}, to: {name: 'ChiefArchitectJD'}, name: 'Chief Architect JD'},
        ]
      }
    },
    org2ImageMap() {
      return {
        height: 200,
        width: 1000,
        hotspots: [
          {coords: {x: 250, y: 1, w: 220, h: 60}, to: {name: 'ChiefArchitectJD'}, name: 'Chief Architect JD'},
        ]
      }
    }
  }
}

</script>

<style scoped>

.orgChart {
  fill: none;
  stroke-width: 2;
  stroke: black;
  width: 1000px;
  height: 200px;
}

.orgText {
  fill: black;
  font-weight: normal;
  stroke: none;
}

.level1 {
  font-size: 1.5rem;
}

.level2 {
  font-size: 0.8rem;
}

.level3 {
  font-size: 0.7rem;
}

</style>