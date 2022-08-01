<template>
  <div class="panel-container">
    <panel-menu :model="topics" class="navigation"/>
    <divider layout="vertical"/>
    <router-view/>
  </div>
</template>

<script>
export default {
  name: "ArchitectureAcademy",
  data() {
    return {
      topics: [
        { label:'What is Architecture?', to: {name: 'ArchitectureIntro' }},
        { label:'Becoming an Architect', to: { name: 'BecomingAnArchitect'}},
        {
          label: 'ARSCIV-X', to: {name: 'ARSCIV_X_Overview'}, items: [
            {label: 'Overview', to: {name: 'ARSCIV_X_Overview'}},
            {label: 'Roles', to: {name: 'ARSCIV_X_Roles'}},
            {label: 'Discussion', to: {name: 'ARSCIV_X_Discussion'}},
            {label: 'Examples', to: {name: 'ARSCIV_X_Examples'}},
          ]
        },
      ]
    }
  }

}
</script>

<style scoped>

</style>