<template>
  <menu-bar :model="items">
    <template #start>
      <img style="padding-right: 10px" alt="logo" :src="'/images/EA Framework Abstract.svg'" height="40" class="mr-2">
    </template>
    <template #end>
      <menu-bar :model="about"/>
    </template>
  </menu-bar>
  <router-view></router-view>
</template>

<script>


export default {
  name: 'App',
  data() {
    return {
      items: [
        { label:'Home', to: {name: 'Home' }},
        { label:'Academy', to: {name: 'AcademyOverview' }},
        { label:'Framework', to: {name: 'FrameworkIntroduction'}},
        { label:'Library', to: {name: 'LibraryOverview'}}
      ],
      about: [
        { icon: 'pi pi-fw pi-arrow-left', command: () => {this.$router.go(-1)} },
        { icon: 'pi pi-fw pi-arrow-right', command: () => {this.$router.go(1)} },
        { label:'About', to: {name: 'About'}},
      ]
    }
  },
}
</script>

<style>

:root {
  --font-family: Avenir, Helvetica, Arial, sans-serif;
  --image-map-hover: lightgray;
  --image-map-on: lightgreen;
}


#app {
  font-family: var(--font-family);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: left;
  color: #2c3e50;
  margin-top: 60px;
}


.text-container {
  max-width: 1000px;
  display: block;
}

.panel-container {
  display: flex;
  padding: 15px 15px 15px 15px;
}

.p-tieredmenu ul {
  white-space: nowrap;
}

.p-panelmenu .p-menuitem-text {
  white-space: nowrap;
}

.caption {
  display: block;
  font-size: 0.8rem;
  text-align: center;
}


.jd {
  display: block;
  margin-left: 40px;
  margin-right: 40px;
  padding: 10px 10px 10px 10px;
  border: 1px solid black;
  font-family: Courier, sans-serif;
  font-size: 0.9rem;
}

.banner {
  display: block;
  margin: 10px 2px 10px 2px;
  text-align: center;
  font-size: 1.2rem;
  background-color: var(--surface-b)
}

.note {
  display: block;
  margin-left: 100px;
  margin-right: 100px;
  padding: 10px 10px 10px 10px;
  border: 1px solid black;
  font-style: italic;
  font-size: 1rem;
}


html {
  height: 100%;
  width: 100%;
}

body {
  margin: 0;
  height: 100%;
  width: 100%;
  overflow-x: hidden;
  overflow-y: auto;
  background-color: var(--surface-a);
  font-family: var(--font-family);
  font-weight: 400;
  color: var(--text-color);
}

</style>
