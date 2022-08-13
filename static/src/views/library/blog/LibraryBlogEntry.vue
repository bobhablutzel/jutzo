<template>
  <div class="text-container">
    <h1>{{blog_entry.title}}</h1>
    <div v-for="(section, index) in blog_entry.body" :key="index" >
      <div v-if="section.header" v-html="getHeaderSection(section)"/>
      <div v-if="section.text" v-html="getTextSection(section)"/>
      <div v-if="section.img" v-html="getImgSection(section)"/>
      <div v-if="section.imageMap">
        <image-map :src="section.imageMap.src" :image-map="section.imageMap.map"/>
      </div>
    </div>
  </div>
</template>

<script>
import ImageMap from "@/components/ImageMap"
import {retrieveBlogEntry} from "@/services/blogService";

export default {
  name: "LibraryBlogEntry",
  components: {ImageMap},
  data() {
    return {
      blog_entry: { title: "Loading…", body: []},
    }
  },
  mounted() {
    retrieveBlogEntry(this.$route.params.id, (entry) => this.blog_entry = entry )
  },
  beforeRouteUpdate(to) {
    retrieveBlogEntry(to.params.id, (entry) => this.blog_entry = entry)
  },
  methods: {
   isInt(value) {
      // noinspection EqualityComparisonWithCoercionJS
     return value && !isNaN(value) &&
        parseInt(value) == value &&
        !isNaN(parseInt(value, 10));
    },
    getHeaderSection(section) {
      const level = this.isInt(section.level) ? parseInt(section.level) : 1
      return "<h" + level.toString() + ">" + section.header + "</h" + level.toString() + ">";
    },
    getTextSection(section) {
      return '<p>' + section.text + '</p>'
    },
    getImgSection(section) {
      return '<img src="' + section.img.src + '" alt="' + (section.img.alt ? section.img.alt : '') + '"/>'
    },
  },
}
</script>


<style scoped>

</style>